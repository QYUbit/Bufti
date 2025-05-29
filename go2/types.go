package bufti2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type BuftiType interface {
	getName() string
	encode(*bytes.Buffer, any) error
	decode(*bytes.Buffer, reflect.Value) error
}

type SimpleType int

const (
	Int8 SimpleType = iota
	Int16
	Int32
	Int64
	Float32
	Float64
	Bool
	String
)

func (t SimpleType) getName() string {
	return "simple"
}

func (t SimpleType) String() string {
	typeNames := []string{"int8", "int16", "int32", "int64", "float64", "float32", "boolean", "string"}
	return fmt.Sprintf("bufti %s", typeNames[t])
}

func (t SimpleType) encode(buf *bytes.Buffer, value any) error {
	switch t {
	case Int8:
		var v int8
		switch val := value.(type) {
		case int8:
			v = val
		case int:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("int value %d out of range for int8", val)
			}
			v = int8(val)
		case int16:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("int16 value %d out of range for int8", val)
			}
			v = int8(val)
		case int32:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("int32 value %d out of range for int8", val)
			}
			v = int8(val)
		case int64:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("int64 value %d out of range for int8", val)
			}
			v = int8(val)
		default:
			return fmt.Errorf("cannot convert %T to int8", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Int16:
		var v int16
		switch val := value.(type) {
		case int16:
			v = val
		case int8:
			v = int16(val)
		case int:
			if val < math.MinInt16 || val > math.MaxInt16 {
				return fmt.Errorf("int value %d out of range for int16", val)
			}
			v = int16(val)
		case int32:
			if val < math.MinInt16 || val > math.MaxInt16 {
				return fmt.Errorf("int32 value %d out of range for int16", val)
			}
			v = int16(val)
		case int64:
			if val < math.MinInt16 || val > math.MaxInt16 {
				return fmt.Errorf("int64 value %d out of range for int16", val)
			}
			v = int16(val)
		default:
			return fmt.Errorf("cannot convert %T to int16", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Int32:
		var v int32
		switch val := value.(type) {
		case int32:
			v = val
		case int8:
			v = int32(val)
		case int16:
			v = int32(val)
		case int:
			if val < math.MinInt32 || val > math.MaxInt32 {
				return fmt.Errorf("int value %d out of range for int32", val)
			}
			v = int32(val)
		case int64:
			if val < math.MinInt32 || val > math.MaxInt32 {
				return fmt.Errorf("int64 value %d out of range for int32", val)
			}
			v = int32(val)
		default:
			return fmt.Errorf("cannot convert %T to int32", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Int64:
		var v int64
		switch val := value.(type) {
		case int64:
			v = val
		case int8:
			v = int64(val)
		case int16:
			v = int64(val)
		case int32:
			v = int64(val)
		case int:
			v = int64(val)
		default:
			return fmt.Errorf("cannot convert %T to int64", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Float32:
		var v float32
		switch val := value.(type) {
		case float32:
			v = val
		case float64:
			if val < -math.MaxFloat32 || val > math.MaxFloat32 {
				return fmt.Errorf("float64 value %f out of range for float32", val)
			}
			v = float32(val)
		default:
			return fmt.Errorf("cannot convert %T to float32", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Float64:
		var v float64
		switch val := value.(type) {
		case float64:
			v = val
		case float32:
			v = float64(val)
		default:
			return fmt.Errorf("cannot convert %T to float64", value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Bool:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("expected bool, got %T", value)
		}
		var b byte
		if v {
			b = 1
		}
		return buf.WriteByte(b)

	case String:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
		if err := binary.Write(buf, binary.LittleEndian, uint32(len(v))); err != nil {
			return err
		}
		_, err := buf.WriteString(v)
		return err

	default:
		return fmt.Errorf("unknown SimpleType: %d", t)
	}
}

func (t SimpleType) decode(buf *bytes.Buffer, val reflect.Value) error {
	switch t {
	case Int8:
		var v int8
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode int8: %w", err)
		}
		return nil

	case Int16:
		var v int16
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode int16: %w", err)
		}
		return nil

	case Int32:
		var v int32
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode int32: %w", err)
		}
		return nil

	case Int64:
		var v int64
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode int64: %w", err)
		}
		return nil

	case Float32:
		var v float32
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode float32: %w", err)
		}
		return nil

	case Float64:
		var v float64
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("failed to decode float64: %w", err)
		}
		return nil

	case Bool:
		_, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to decode bool: %w", err)
		}
		return nil

	case String:
		var length uint32
		if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
			return fmt.Errorf("failed to decode string length: %w", err)
		}

		if length > uint32(buf.Len()) {
			return fmt.Errorf("string length %d exceeds buffer size %d", length, buf.Len())
		}

		data := make([]byte, length)
		n, err := buf.Read(data)
		if err != nil {
			return fmt.Errorf("failed to decode string data: %w", err)
		}
		if n != int(length) {
			return fmt.Errorf("expected to read %d bytes, got %d", length, n)
		}

		return nil

	default:
		return fmt.Errorf("unknown SimpleType: %d", t)
	}
}

type ListType struct {
	elementType BuftiType
}

func List(elementType BuftiType) ListType {
	return ListType{elementType: elementType}
}

func (t ListType) getName() string {
	return "list"
}

func (t ListType) String() string {
	return fmt.Sprintf("bufti list of %ss", t.elementType)
}

func (t ListType) encode(buf *bytes.Buffer, value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return fmt.Errorf("%w: can not apply value of type %T to %s", fmt.Errorf(""), value, "")
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(val.Len())); err != nil {
		return err
	}

	for i := range val.Len() {
		if !val.CanInterface() {
			continue
		}
		if err := t.elementType.encode(buf, val.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (t ListType) decode(buf *bytes.Buffer, val reflect.Value) error {
	return nil
}

type MapType struct {
	keyType   SimpleType
	valueType BuftiType
}

func Map(keyType SimpleType, valueType BuftiType) MapType {
	return MapType{keyType: keyType, valueType: valueType}
}

func (t MapType) getName() string {
	return "map"
}

func (t MapType) String() string {
	return fmt.Sprintf("bufti map from %s to %s", t.keyType, t.valueType)
}

func (t MapType) encode(buf *bytes.Buffer, value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Map {
		return fmt.Errorf("%w: can not apply value of type %T to %s", fmt.Errorf(""), value, "")
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(val.Len())); err != nil {
		return err
	}

	for _, key := range val.MapKeys() {
		if !key.CanInterface() || val.MapIndex(key).CanInterface() {
			continue
		}
		if err := t.keyType.encode(buf, key.Interface()); err != nil {
			return err
		}
		if err := t.valueType.encode(buf, val.MapIndex(key).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (t MapType) decode(buf *bytes.Buffer, val reflect.Value) error {
	return nil
}

type ReferenceType struct {
	model *Model
}

func Reference(model *Model) ReferenceType {
	return ReferenceType{model: model}
}

func (t ReferenceType) getName() string {
	return "reference"
}

func (t ReferenceType) String() string {
	return fmt.Sprintf("bufti model %s", t.model.name)
}

func (t ReferenceType) encode(buf *bytes.Buffer, data any) error {
	return t.model.encode(buf, data)
}

func (t ReferenceType) decode(buf *bytes.Buffer, val reflect.Value) error {
	return nil
}
