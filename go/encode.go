package butil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Encode serializes the given data according to the model schema.
// The data can be a struct or map[string]any. Struct fields are mapped to
// schema fields using either the field name or the `butil` tag.
//
// Returns ErrInput if the data is nil or of an unsupported type.
// Returns ErrModel if required fields are missing or schema validation fails.
func (m *Model) Encode(data any) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("%w: cannot encode nil", ErrInput)
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	if err := binary.Write(buf, binary.LittleEndian, ProtocolVersion); err != nil {
		return nil, fmt.Errorf("failed to write protocol version")
	}

	if err := m.encode(buf, reflect.TypeOf(data), reflect.ValueOf(data)); err != nil {
		return nil, err
	}

	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

func (m *Model) encode(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
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
		return fmt.Errorf("%w: invalid input type %v", ErrInput, t.Kind())
	}
	return nil
}

type valueFieldPair struct {
	v     reflect.Value
	field ModelField
}

func (m *Model) encodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
	m.mu.RLock()
	fieldMap, exists := m.fieldCache[t]
	m.mu.RUnlock()

	if !exists {
		fieldMap = make(map[string]reflect.Value, len(m.schema))

		for i := 0; i < v.NumField(); i++ {
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

		m.mu.Lock()
		if m.fieldCache == nil {
			m.fieldCache = make(map[reflect.Type]map[string]reflect.Value)
		}
		m.fieldCache[t] = fieldMap
		m.mu.Unlock()
	}

	valueFieldPairs := make(map[byte]valueFieldPair, len(m.schema))
	for fieldName, value := range fieldMap {
		index, exists := m.labels[fieldName]
		if !exists {
			return fmt.Errorf("%w: field %s not found in model %s", ErrInput, fieldName, m.name)
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d not found on model %s", ErrModel, index, m.name)
		}

		valueFieldPairs[index] = valueFieldPair{v: value, field: schemaField}
	}

	for index, field := range m.schema {
		if field.isRequired == nil {
			continue
		}
		if !*field.isRequired {
			continue
		}
		if _, exists := valueFieldPairs[index]; !exists {
			return fmt.Errorf("%w: required field %s is missing for model %s", ErrInput, field.label, m.name)
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, uint32(len(valueFieldPairs))); err != nil {
		return err
	}

	for index, pair := range valueFieldPairs {
		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := pair.field.fieldType.Encode(buf, pair.v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) encodeMap(buf *bytes.Buffer, _ reflect.Type, v reflect.Value) error {
	fieldMap := make(map[byte]valueFieldPair, len(m.schema))

	for _, k := range v.MapKeys() {
		key := k.Interface()
		value := v.MapIndex(k)

		strKey, ok := key.(string)
		if !ok {
			return fmt.Errorf("%w: map key has to be a string, intead: %T", ErrInput, key)
		}

		index, exists := m.labels[strKey]
		if !exists {
			return fmt.Errorf("%w: field %s not found in model %s", ErrInput, strKey, m.name)
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index %d not found on model %s", ErrModel, index, m.name)
		}

		fieldMap[index] = valueFieldPair{v: value, field: schemaField}
	}

	for index, field := range m.schema {
		if field.isRequired == nil {
			continue
		}
		if !*field.isRequired {
			continue
		}
		if _, exists := fieldMap[index]; !exists {
			return fmt.Errorf("%w: required field %s is missing for model %s", ErrInput, field.label, m.name)
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, uint32(v.Len())); err != nil {
		return err
	}

	for index, pair := range fieldMap {
		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := pair.field.fieldType.Encode(buf, pair.v); err != nil {
			return err
		}
	}
	return nil
}
