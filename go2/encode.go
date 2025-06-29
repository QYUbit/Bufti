package bufti2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func (m *Model) Encode(data any) ([]byte, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()

	if err := binary.Write(buf, binary.LittleEndian, ProtocolVersion); err != nil {
		return nil, fmt.Errorf("failed to write protocol version")
	}

	if err := m.encode(buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Model) encode(buf *bytes.Buffer, data any) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	v = indirectValue(v)
	t = indirectType(t)

	switch t.Kind() {
	case reflect.Struct:
		if err := m.encodeStruct(buf, t, v); err != nil {
			return err
		}
	case reflect.Map:
		if err := m.encodeMap(buf, t, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: invalid input type %T", ErrInput, data)
	}
	return nil
}

type valueFieldPair struct {
	v     reflect.Value
	field ModelField
}

func (m *Model) encodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
	fieldMap := make(map[byte]valueFieldPair)

	for i := range v.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanInterface() {
			continue
		}

		fieldName := field.Name
		if tag := field.Tag.Get("bufti"); tag != "" {
			fieldName = tag
		}

		index, exists := m.labels[fieldName]
		if !exists {
			continue
		}

		schemaField, exists := m.schema[index]
		if !exists && *(schemaField.isRequired) {
			return fmt.Errorf("%w: index %d not found on model %s", ErrModel, index, m.name)
		}
		if !exists {
			continue
		}

		fieldMap[index] = valueFieldPair{v: value, field: schemaField}
	}

	if err := binary.Write(buf, binary.LittleEndian, uint32(len(fieldMap))); err != nil {
		return err
	}

	for index, pair := range fieldMap {
		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := pair.field.fieldType.Encode(buf, pair.v.Interface()); err != nil { // ! Why interface() here and don't just keep the reflect value?
			return err
		}
	}
	return nil
}

func (m *Model) encodeMap(buf *bytes.Buffer, _ reflect.Type, v reflect.Value) error {
	fieldMap := make(map[byte]valueFieldPair)

	for _, k := range v.MapKeys() {
		key := k.Interface()
		value := v.MapIndex(k)

		strKey, ok := key.(string)
		if !ok {
			return fmt.Errorf("%w: map key has to be a string, intead: %T", ErrInput, key)
		}

		index, exists := m.labels[strKey]
		if !exists {
			continue
		}

		schemaField, exists := m.schema[index]
		if !exists && *(schemaField.isRequired) {
			return fmt.Errorf("%w: index %d not found on model %s", ErrModel, index, m.name)
		}
		if !exists {
			continue
		}

		fieldMap[index] = valueFieldPair{v: value, field: schemaField}
	}

	if err := binary.Write(buf, binary.LittleEndian, uint32(v.Len())); err != nil {
		return err
	}

	for index, pair := range fieldMap {
		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := pair.field.fieldType.Encode(buf, pair.v.Interface()); err != nil {
			return err
		}
	}
	return nil
}
