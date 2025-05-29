package bufti2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type Decoder struct {
	buf      *bytes.Buffer
	depth    int
	maxDepth int
}

func NewDecoder(data []byte) *Decoder {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Write(data)

	return &Decoder{
		buf:      buf,
		maxDepth: 100,
	}
}

func (d *Decoder) Close() {
	if d.buf != nil {
		bufferPool.Put(d.buf)
		d.buf = nil
	}
}

// Decode in beliebige Zielstruktur
func (d *Decoder) Decode(target any) error {
	if d.depth > d.maxDepth {
		return fmt.Errorf("maximum recursion depth exceeded")
	}

	var version uint32
	if err := binary.Read(d.buf, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != ProtocolVersion {
		return errors.New("version error")
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := targetValue.Elem()
	return d.decodeValue(elem)
}

func (d *Decoder) decodeValue(v reflect.Value) error {
	d.depth++
	defer func() { d.depth-- }()

	switch v.Kind() {
	case reflect.Int:
		val, err := d.decodeInt64()
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case reflect.Int8:
		val, err := d.decodeInt8()
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case reflect.Int16:
		val, err := d.decodeInt16()
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case reflect.Int32:
		val, err := d.decodeInt32()
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case reflect.Int64:
		val, err := d.decodeInt64()
		if err != nil {
			return err
		}
		v.SetInt(val)

	case reflect.Float32:
		val, err := d.decodeFloat32()
		if err != nil {
			return err
		}
		v.SetFloat(float64(val))

	case reflect.Float64:
		val, err := d.decodeFloat64()
		if err != nil {
			return err
		}
		v.SetFloat(val)

	case reflect.Bool:
		val, err := d.decodeBool()
		if err != nil {
			return err
		}
		v.SetBool(val)

	case reflect.String:
		val, err := d.decodeString()
		if err != nil {
			return err
		}
		v.SetString(val)

	case reflect.Slice:
		return d.decodeSlice(v)

	case reflect.Map:
		return d.decodeMap(v)

	case reflect.Struct:
		return d.decodeStruct(v)

	case reflect.Ptr:
		return d.decodePointer(v)

	default:
		return fmt.Errorf("unsupported type: %v", v.Kind())
	}

	return nil
}

// Primitive Decoder
func (d *Decoder) decodeInt8() (int8, error) {
	var v int8
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeInt16() (int16, error) {
	var v int16
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeInt32() (int32, error) {
	var v int32
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeInt64() (int64, error) {
	var v int64
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeFloat32() (float32, error) {
	var v float32
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeFloat64() (float64, error) {
	var v float64
	err := binary.Read(d.buf, binary.BigEndian, &v)
	return v, err
}

func (d *Decoder) decodeBool() (bool, error) {
	b, err := d.buf.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func (d *Decoder) decodeString() (string, error) {
	var length uint32
	if err := binary.Read(d.buf, binary.BigEndian, &length); err != nil {
		return "", err
	}

	if length == 0 {
		return "", nil
	}

	data := make([]byte, length)
	n, err := d.buf.Read(data)
	if err != nil {
		return "", err
	}
	if n != int(length) {
		return "", fmt.Errorf("expected %d bytes, got %d", length, n)
	}

	return string(data), nil
}

// Komplexe Strukturen
func (d *Decoder) decodeSlice(v reflect.Value) error {
	// Slice-Länge lesen
	var length uint32
	if err := binary.Read(d.buf, binary.BigEndian, &length); err != nil {
		return err
	}

	// Slice erstellen
	_ = v.Type().Elem()
	slice := reflect.MakeSlice(v.Type(), int(length), int(length))

	// Elemente decodieren
	for i := range int(length) {
		elem := slice.Index(i)
		if err := d.decodeValue(elem); err != nil {
			return fmt.Errorf("error decoding slice element %d: %w", i, err)
		}
	}

	v.Set(slice)
	return nil
}

func (d *Decoder) decodeMap(v reflect.Value) error {
	// Map-Länge lesen
	var length uint32
	if err := binary.Read(d.buf, binary.BigEndian, &length); err != nil {
		return err
	}

	// Map erstellen
	mapType := v.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()
	newMap := reflect.MakeMap(mapType)

	// Key-Value Paare decodieren
	for i := 0; i < int(length); i++ {
		// Key decodieren
		key := reflect.New(keyType).Elem()
		if err := d.decodeValue(key); err != nil {
			return fmt.Errorf("error decoding map key %d: %w", i, err)
		}

		// Value decodieren
		value := reflect.New(valueType).Elem()
		if err := d.decodeValue(value); err != nil {
			return fmt.Errorf("error decoding map value %d: %w", i, err)
		}

		newMap.SetMapIndex(key, value)
	}

	v.Set(newMap)
	return nil
}

func (d *Decoder) decodeStruct(v reflect.Value) error {
	structType := v.Type()

	for i := range structType.NumField() {
		field := v.Field(i)
		structField := structType.Field(i)

		// Exportierte Felder überspringen
		if !field.CanSet() {
			continue
		}

		// Tags für Steuerung verwenden (optional)
		tag := structField.Tag.Get("binary")
		if tag == "-" {
			continue // Feld überspringen
		}

		if err := d.decodeValue(field); err != nil {
			return fmt.Errorf("error decoding struct field %s: %w", structField.Name, err)
		}
	}

	return nil
}

func (d *Decoder) decodePointer(v reflect.Value) error {
	// Null-Pointer Check (erstes Byte)
	isNull, err := d.decodeBool()
	if err != nil {
		return err
	}

	if isNull {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	// Neuen Wert erstellen und decodieren
	elem := reflect.New(v.Type().Elem())
	if err := d.decodeValue(elem.Elem()); err != nil {
		return err
	}

	v.Set(elem)
	return nil
}
