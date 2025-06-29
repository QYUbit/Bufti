package bufti2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// TODO Indirect values before decoding

type BuftiType interface {
	Encode(*bytes.Buffer, reflect.Value) error
	Decode(*bytes.Buffer, reflect.Value) error
}

type SimpleType int

const (
	Bool SimpleType = iota
	Uint8
	Uint16
	Uint32
	Uint64
	Int8
	Int16
	Int32
	Int64
	Float32
	Float64
	Bytes
	String
)

func (t SimpleType) String() string {
	typeNames := [13]string{"boolean", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "bytes", "string"}
	return fmt.Sprintf("bufti %s", typeNames[t])
}

func (t SimpleType) reflectType() (reflect.Type, error) {
	switch t {
	case Int8:
		return reflect.TypeOf(int8(0)), nil
	case Int16:
		return reflect.TypeOf(int16(0)), nil
	case Int32:
		return reflect.TypeOf(int32(0)), nil
	case Int64:
		return reflect.TypeOf(int64(0)), nil
	case Float32:
		return reflect.TypeOf(float32(0)), nil
	case Float64:
		return reflect.TypeOf(float64(0)), nil
	case Bool:
		return reflect.TypeOf(bool(false)), nil
	case String:
		return reflect.TypeOf(string("")), nil
	default:
		return nil, fmt.Errorf("%v is no simple type", t)
	}
}

type ListType struct {
	elementType BuftiType
}

func List(elementType BuftiType) ListType {
	return ListType{elementType: elementType}
}

func (t ListType) String() string {
	return fmt.Sprintf("bufti list of %ss", t.elementType)
}

func (t ListType) Encode(buf *bytes.Buffer, val reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("can not encode value of type %v as %s", val.Kind(), t)
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(val.Len())); err != nil {
		return err
	}

	for i := range val.Len() {
		if !val.CanInterface() {
			continue
		}
		if err := t.elementType.Encode(buf, val.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (t ListType) Decode(buf *bytes.Buffer, v reflect.Value) error {
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return fmt.Errorf("%w: failed to decode list length: %w", ErrBuffer, err)
	}

	// TODO Make seperate interface functions for each decode to reduce reflect calls
	var slice reflect.Value
	if v.Kind() == reflect.Slice {
		slice = reflect.MakeSlice(v.Type(), int(length), int(length))
	} else {
		slice = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((*any)(nil)).Elem()), int(length), int(length))
	}

	for i := range int(length) {
		elem := slice.Index(i)
		if err := t.elementType.Decode(buf, elem); err != nil {
			return err
		}
	}

	v.Set(slice)
	return nil
}

type MapType struct {
	keyType   SimpleType
	valueType BuftiType
}

func Map(keyType SimpleType, valueType BuftiType) MapType {
	return MapType{keyType: keyType, valueType: valueType}
}

func (t MapType) String() string {
	return fmt.Sprintf("bufti map (%s -> %s)", t.keyType, t.valueType)
}

func (t MapType) Encode(buf *bytes.Buffer, val reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	if val.Kind() != reflect.Map {
		return fmt.Errorf("can not encode value of type %v as %s", val.Kind(), t)
	}
	if err := binary.Write(buf, binary.LittleEndian, uint32(val.Len())); err != nil {
		return err
	}

	for _, key := range val.MapKeys() {
		if !key.CanInterface() || !val.MapIndex(key).CanInterface() {
			continue
		}
		if err := t.keyType.Encode(buf, key); err != nil {
			return err
		}
		if err := t.valueType.Encode(buf, val.MapIndex(key)); err != nil {
			return err
		}
	}
	return nil
}

func (t MapType) Decode(buf *bytes.Buffer, v reflect.Value) error {
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return fmt.Errorf("%w: failed to decode map length: %w", ErrBuffer, err)
	}

	var newMap reflect.Value
	if v.Kind() == reflect.Interface {
		keyType, err := t.keyType.reflectType()
		if err != nil {
			return err
		}

		mapType := reflect.MapOf(keyType, reflect.TypeOf((*any)(nil)).Elem())
		newMap = reflect.MakeMap(mapType)
	} else {
		newMap = reflect.MakeMap(v.Type())
	}

	for range length {
		keyValue := reflect.New(newMap.Type().Key()).Elem()
		if err := t.keyType.Decode(buf, keyValue); err != nil {
			return err
		}

		valueValue := reflect.New(newMap.Type().Elem()).Elem()
		if err := t.valueType.Decode(buf, valueValue); err != nil {
			return err
		}

		newMap.SetMapIndex(keyValue, valueValue)
	}

	v.Set(newMap)
	return nil
}

type ReferenceType struct {
	model *Model
}

func Reference(model *Model) ReferenceType {
	return ReferenceType{model: model}
}

func (t ReferenceType) String() string {
	return fmt.Sprintf("bufti model %s", t.model.name)
}

func (t ReferenceType) Encode(buf *bytes.Buffer, val reflect.Value) error {
	return t.model.encode(buf, val.Type(), val)
}

func (t ReferenceType) Decode(buf *bytes.Buffer, val reflect.Value) error {
	return t.model.decode(buf, val.Type(), val)
}

func indirectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
