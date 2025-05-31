package bufti2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

func (m *Model) Decode(data []byte, dest any) error {
	v := reflect.ValueOf(dest)
	t := reflect.TypeOf(dest)

	if t.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("dest has to be a pointer")
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
		return err
	}
	if version != ProtocolVersion {
		return fmt.Errorf("bufti: incompatible bufti version: this package uses version %d input uses version %d", ProtocolVersion, version)
	}

	return m.decode(buf, t, v, len(data)/2)
}

func (m *Model) decode(buf *bytes.Buffer, t reflect.Type, v reflect.Value, maxIterations int) error {
	switch v.Kind() {
	case reflect.Struct:
		return m.decodeStruct(buf, t, v, maxIterations)
	case reflect.Map:
		if t.Key().Kind() != reflect.String || t.Elem().Kind() != reflect.Interface {
			return fmt.Errorf("bufti: destination has to be a string to any map, instead: %s", t.String())
		}
		return m.decodeMap(buf, t, v, maxIterations)
	default:
		return fmt.Errorf("bufti: invalid destination type %s", t.String())
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
			return fmt.Errorf("bufti: index %d does not exist on this model", index)
		}

		if err = schemaField.fieldType.decode(buf, fieldMap[schemaField.label]); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) decodeMap(buf *bytes.Buffer, _ reflect.Type, v reflect.Value, maxIterations int) error {
	//decodeMap := reflect.MakeMap(t)

	for range maxIterations {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("bufti: index %d does not exist on this model", index)
		}

		//if !slices.Contains(v.MapKeys(), reflect.ValueOf(schemaField.label)) {
		//	fmt.Println(v.MapKeys())
		//	return fmt.Errorf("bufti could not find label %s in map", schemaField.label)
		//}
		var mapValue any
		if err = schemaField.fieldType.decode(buf, reflect.ValueOf(&mapValue).Elem()); err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(schemaField.label), reflect.ValueOf(mapValue))
	}

	//v.Set(decodeMap)
	return nil
}
