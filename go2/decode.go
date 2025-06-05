package bufti2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func (m *Model) Decode(data []byte, dest any) error {
	v := reflect.ValueOf(dest)
	t := reflect.TypeOf(dest)

	if t.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("%w: dest has to be a pointer, instead: %T", ErrInput, dest)
	}
	v = v.Elem()
	t = t.Elem()

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()
	if _, err := buf.Write(data); err != nil {
		return err
	}

	var version uint32
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("failed to read protocol version")
	}
	if version != ProtocolVersion {
		return fmt.Errorf("incompatible bufti version: this package uses version %d, buffer uses version %d", ProtocolVersion, version)
	}

	return m.decode(buf, t, v, len(data)/2)
}

func (m *Model) decode(buf *bytes.Buffer, t reflect.Type, v reflect.Value, maxIterations int) error {
	switch v.Kind() {
	case reflect.Struct:
		return m.decodeStruct(buf, t, v, maxIterations)
	case reflect.Map:
		return m.decodeMap(buf, t, v, maxIterations)
	default:
		return fmt.Errorf("%w: invalid destination type %s", ErrInput, t.String())
	}
}

func (m *Model) decodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value, maxIterations int) error {
	fieldMap := make(map[string]reflect.Value)

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		fieldName := field.Name
		if tag := field.Tag.Get("bufti"); tag != "" {
			fieldName = tag
		}

		fieldMap[fieldName] = value
	}

	for range maxIterations {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d does not exist on model %s", ErrModel, index, m.name)
		}

		if err = schemaField.fieldType.Decode(buf, fieldMap[schemaField.label]); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) decodeMap(buf *bytes.Buffer, t reflect.Type, v reflect.Value, maxIterations int) error {
	if t.Key().Kind() != reflect.String || t.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("%w: destination has to be a map[string]any, instead: %T", ErrInput, t.String())
	}

	for range maxIterations {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d does not exist on model %s", ErrModel, index, m.name)
		}

		var mapValue any
		if err = schemaField.fieldType.Decode(buf, reflect.ValueOf(&mapValue).Elem()); err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(schemaField.label), reflect.ValueOf(mapValue))
	}
	return nil
}
