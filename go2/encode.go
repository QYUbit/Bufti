package bufti2

import (
	"bytes"
	"fmt"
	"reflect"
)

func (m *Model) Encode(data any) ([]byte, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

	buf.Reset()

	if err := m.encode(buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Model) encode(buf *bytes.Buffer, data any) error {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

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
		// Error
	}
	return nil
}

func (m *Model) encodeStruct(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
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
		fmt.Printf("field %s name %s exists %t\n", field.Name, fieldName, exists)
		if !exists {
			continue
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index not found (%d)", fmt.Errorf(""), index)
		}

		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := schemaField.fieldType.encode(buf, value.Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) encodeMap(buf *bytes.Buffer, _ reflect.Type, v reflect.Value) error {
	if v.Kind() != reflect.Map {
		return fmt.Errorf("")
	}

	for _, k := range v.MapKeys() {
		key := k.Interface()
		value := v.MapIndex(k).Interface()

		strKey, ok := key.(string)
		if !ok {
			return fmt.Errorf("")
		}

		index, exists := m.labels[strKey]
		if !exists {
			continue
		}

		schemaField, exists := m.schema[index]
		if !exists {
			return fmt.Errorf("%w: index not found (%d)", fmt.Errorf(""), index)
		}

		if err := buf.WriteByte(index); err != nil {
			return err
		}

		if err := schemaField.fieldType.encode(buf, value); err != nil {
			return err
		}
	}
	return nil
}
