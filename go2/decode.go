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

	case reflect.Map:
		if t.Key().Kind() != reflect.String || t.Elem().Kind() != reflect.Interface {
			return fmt.Errorf("bufti: destination has to be a string to any map, instead: %s", t.String())
		}
	default:
		return fmt.Errorf("bufti: invalid destination type %s", t.String())
	}

	for range maxIterations {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("bufto: index %d does not exist on this model", index)
		}

		if err = schemaField.fieldType.decode(buf, v); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) decodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value, maxIterations int) error {
	for range maxIterations {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("bufto: index %d does not exist on this model", index)
		}

		//t.Field()

		if err = schemaField.fieldType.decode(buf, v); err != nil {
			return err
		}
	}
	return nil
}
