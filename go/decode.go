package butil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Decode deserializes binary data into the given destination according to the model schema.
// The destination must be a pointer to a struct or map[string]any.
// Struct fields are mapped from schema fields using either the field name or the `butil` tag.
//
// Returns ErrInput if dest is not a pointer or is nil.
// Returns ErrVersion if the data was encoded with an incompatible protocol version.
// Returns ErrBuffer if the data is corrupted or cannot be parsed.
// Returns ErrModel if the data references fields not defined in the schema.
func (m *Model) Decode(data []byte, dest any) error {
	v := reflect.ValueOf(dest)
	t := reflect.TypeOf(dest)

	if t.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("%w: dest has to be a pointer, instead: %s", ErrInput, t.Kind())
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()
	if _, err := buf.Write(data); err != nil {
		return err
	}

	var version uint32
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("%w: failed to read protocol version", ErrBuffer)
	}
	if version != ProtocolVersion {
		return fmt.Errorf("%w: incompatible butil version: this package uses version %d, buffer uses version %d", ErrVersion, ProtocolVersion, version)
	}

	return m.decode(buf, t, v)
}

func (m *Model) decode(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
	var fieldCount uint32
	if err := binary.Read(buf, binary.LittleEndian, &fieldCount); err != nil {
		return err
	}

	v = indirectValue(v)
	t = indirectType(t)

	switch t.Kind() {
	case reflect.Struct:
		return m.decodeStruct(buf, t, v, int(fieldCount))
	case reflect.Map:
		return m.decodeMap(buf, t, v, int(fieldCount))
	default:
		return fmt.Errorf("%w: invalid destination type %s", ErrInput, t.Kind())
	}
}

func (m *Model) decodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value, fieldCount int) error {
	fieldMap := make(map[string]reflect.Value, len(m.schema))

	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanInterface() {
			continue
		}

		fieldName := field.Name
		if tag := field.Tag.Get("butil"); tag != "" {
			fieldName = tag
		}

		fieldMap[fieldName] = value
	}

	for _, field := range m.schema {
		if field.isRequired == nil {
			continue
		}
		if !*field.isRequired {
			continue
		}
		if _, exists := fieldMap[field.label]; !exists {
			return fmt.Errorf("%w: required field %s is missing for model %s", ErrBuffer, field.label, m.name)
		}
	}

	for range fieldCount {
		index, err := buf.ReadByte()
		if err != nil {
			return err
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d does not exist on model %s", ErrBuffer, index, m.name)
		}

		if err = schemaField.fieldType.Decode(buf, fieldMap[schemaField.label]); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) decodeMap(buf *bytes.Buffer, t reflect.Type, v reflect.Value, fieldCount int) error {
	if t.Key().Kind() != reflect.String || t.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("%w: destination has to be a map[string]any, instead: %T", ErrInput, t.String())
	}

	decodedFields := make(map[byte]bool)
	for range fieldCount {
		index, err := buf.ReadByte()
		if err != nil {
			break
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d does not exist on model %s", ErrBuffer, index, m.name)
		}

		decodedFields[index] = true

		var mapValue any
		if err = schemaField.fieldType.Decode(buf, reflect.ValueOf(&mapValue).Elem()); err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(schemaField.label), reflect.ValueOf(mapValue))
	}

	for index, field := range m.schema {
		if field.isRequired == nil {
			continue
		}
		if !*field.isRequired {
			continue
		}
		if !decodedFields[index] {
			return fmt.Errorf("%w: required field %s is missing for model %s", ErrBuffer, field.label, m.name)
		}
	}
	return nil
}
