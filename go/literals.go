package bufti2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

func (t SimpleType) Encode(buf *bytes.Buffer, reflectValue reflect.Value) error {
	if !reflectValue.CanInterface() {
		return fmt.Errorf("%w: value cannot be converted to a interface interface", ErrInput)
	}
	value := reflectValue.Interface()

	switch t {
	case Int8:
		var v int8
		switch val := value.(type) {
		case int8:
			v = val
		case int:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("%w: int value %d out of range for int8", ErrInput, val)
			}
			v = int8(val)
		case int16:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("%w: int16 value %d out of range for int8", ErrInput, val)
			}
			v = int8(val)
		case int32:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("%w: int32 value %d out of range for int8", ErrInput, val)
			}
			v = int8(val)
		case int64:
			if val < math.MinInt8 || val > math.MaxInt8 {
				return fmt.Errorf("%w: int64 value %d out of range for int8", ErrInput, val)
			}
			v = int8(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to int8", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Uint8:
		var v uint8
		switch val := value.(type) {
		case uint8:
			v = val
		case uint:
			if val > math.MaxUint8 {
				return fmt.Errorf("%w: uint value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case uint16:
			if val > math.MaxUint8 {
				return fmt.Errorf("%w: uint16 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case uint32:
			if val > math.MaxUint8 {
				return fmt.Errorf("%w: uint32 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case uint64:
			if val > math.MaxUint8 {
				return fmt.Errorf("%w: uint64 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case int:
			if val < 0 || val > math.MaxUint8 {
				return fmt.Errorf("%w: int value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case int8:
			if val < 0 {
				return fmt.Errorf("%w: int8 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case int16:
			if val < 0 || val > math.MaxUint8 {
				return fmt.Errorf("%w: int16 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case int32:
			if val < 0 || val > math.MaxUint8 {
				return fmt.Errorf("%w: int32 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		case int64:
			if val < 0 || val > math.MaxUint8 {
				return fmt.Errorf("%w: int64 value %d out of range for uint8", ErrInput, val)
			}
			v = uint8(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to uint8", ErrInput, value)
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
				return fmt.Errorf("%w: int value %d out of range for int16", ErrInput, val)
			}
			v = int16(val)
		case int32:
			if val < math.MinInt16 || val > math.MaxInt16 {
				return fmt.Errorf("%w: int32 value %d out of range for int16", ErrInput, val)
			}
			v = int16(val)
		case int64:
			if val < math.MinInt16 || val > math.MaxInt16 {
				return fmt.Errorf("%w: int64 value %d out of range for int16", ErrInput, val)
			}
			v = int16(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to int16", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Uint16:
		var v uint16
		switch val := value.(type) {
		case uint16:
			v = val
		case uint8:
			v = uint16(val)
		case uint:
			if val > math.MaxUint16 {
				return fmt.Errorf("%w: uint value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case uint32:
			if val > math.MaxUint16 {
				return fmt.Errorf("%w: uint32 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case uint64:
			if val > math.MaxUint16 {
				return fmt.Errorf("%w: uint64 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case int:
			if val < 0 || val > math.MaxUint16 {
				return fmt.Errorf("%w: int value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case int8:
			if val < 0 {
				return fmt.Errorf("%w: int8 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case int16:
			if val < 0 {
				return fmt.Errorf("%w: int16 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case int32:
			if val < 0 || val > math.MaxUint16 {
				return fmt.Errorf("%w: int32 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		case int64:
			if val < 0 || val > math.MaxUint16 {
				return fmt.Errorf("%w: int64 value %d out of range for uint16", ErrInput, val)
			}
			v = uint16(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to uint16", ErrInput, value)
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
				return fmt.Errorf("%w: int value %d out of range for int32", ErrInput, val)
			}
			v = int32(val)
		case int64:
			if val < math.MinInt32 || val > math.MaxInt32 {
				return fmt.Errorf("%w: int64 value %d out of range for int32", ErrInput, val)
			}
			v = int32(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to int32", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Uint32:
		var v uint32
		switch val := value.(type) {
		case uint32:
			v = val
		case uint8:
			v = uint32(val)
		case uint16:
			v = uint32(val)
		case uint:
			if val > math.MaxUint32 {
				return fmt.Errorf("%w: uint value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case uint64:
			if val > math.MaxUint32 {
				return fmt.Errorf("%w: uint64 value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case int:
			if val < 0 || val > math.MaxUint32 {
				return fmt.Errorf("%w: int value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case int8:
			if val < 0 {
				return fmt.Errorf("%w: int8 value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case int16:
			if val < 0 {
				return fmt.Errorf("%w: int16 value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case int32:
			if val < 0 {
				return fmt.Errorf("%w: int32 value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		case int64:
			if val < 0 || val > math.MaxUint32 {
				return fmt.Errorf("%w: int64 value %d out of range for uint32", ErrInput, val)
			}
			v = uint32(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to uint32", ErrInput, value)
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
			return fmt.Errorf("%w: cannot convert %T to int64", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Uint64:
		var v uint64
		switch val := value.(type) {
		case uint64:
			v = val
		case uint8:
			v = uint64(val)
		case uint16:
			v = uint64(val)
		case uint32:
			v = uint64(val)
		case uint:
			v = uint64(val)
		case int:
			if val < 0 {
				return fmt.Errorf("%w: int value %d out of range for uint64", ErrInput, val)
			}
			v = uint64(val)
		case int8:
			if val < 0 {
				return fmt.Errorf("%w: int8 value %d out of range for uint64", ErrInput, val)
			}
			v = uint64(val)
		case int16:
			if val < 0 {
				return fmt.Errorf("%w: int16 value %d out of range for uint64", ErrInput, val)
			}
			v = uint64(val)
		case int32:
			if val < 0 {
				return fmt.Errorf("%w: int32 value %d out of range for uint64", ErrInput, val)
			}
			v = uint64(val)
		case int64:
			if val < 0 {
				return fmt.Errorf("%w: int64 value %d out of range for uint64", ErrInput, val)
			}
			v = uint64(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to uint64", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Float32:
		var v float32
		switch val := value.(type) {
		case float32:
			v = val
		case float64:
			if val < -math.MaxFloat32 || val > math.MaxFloat32 {
				return fmt.Errorf("%w: float64 value %f out of range for float32", ErrInput, val)
			}
			v = float32(val)
		default:
			return fmt.Errorf("%w: cannot convert %T to float32", ErrInput, value)
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
			return fmt.Errorf("%w: cannot convert %T to float64", ErrInput, value)
		}
		return binary.Write(buf, binary.LittleEndian, v)

	case Bool:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%w: expected bool, got %T", ErrInput, value)
		}
		var b byte
		if v {
			b = 1
		}
		return buf.WriteByte(b)

	case String:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("%w: expected string, got %T", ErrInput, value)
		}
		if err := binary.Write(buf, binary.LittleEndian, uint32(len(v))); err != nil {
			return err
		}
		_, err := buf.WriteString(v)
		return err

	case Bytes:
		v, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("%w: expected byte slice, got %T", ErrInput, value)
		}
		if err := binary.Write(buf, binary.LittleEndian, uint32(len(v))); err != nil {
			return err
		}
		_, err := buf.Write(v)
		return err

	default:
		return fmt.Errorf("%w: unknown SimpleType: %d", ErrModel, t)
	}
}

func (t SimpleType) Decode(buf *bytes.Buffer, val reflect.Value) error {
	switch t {
	case Int8:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set int8 value", ErrInput)
		}
		var v int8
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode int8: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Int8:
			val.SetInt(int64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set int8 value to %s", ErrInput, val.Kind())
		}

	case Uint8:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set uint8 value", ErrInput)
		}
		var v uint8
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode uint8: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Uint8:
			val.SetUint(uint64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set uint8 value to %s", ErrInput, val.Kind())
		}

	case Int16:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set int16 value", ErrInput)
		}
		var v int16
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode int16: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Int16:
			val.SetInt(int64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set int16 value to %s", ErrInput, val.Kind())
		}

	case Uint16:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set uint16 value", ErrInput)
		}
		var v uint16
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode uint16: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Uint16:
			val.SetUint(uint64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set uint16 value to %s", ErrInput, val.Kind())
		}

	case Int32:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set int32 value", ErrInput)
		}
		var v int32
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode int32: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Int32:
			val.SetInt(int64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set int32 value to %s", ErrInput, val.Kind())
		}

	case Uint32:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set uint32 value", ErrInput)
		}
		var v uint32
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode uint32: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Uint:
			val.Set(reflect.ValueOf(uint(v)))
		case reflect.Uint32:
			val.SetUint(uint64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set uint32 value to %s", ErrInput, val.Kind())
		}

	case Int64:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set int64 value", ErrInput)
		}
		var v int64
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode int64: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Int:
			val.Set(reflect.ValueOf(int(v)))
		case reflect.Int64:
			val.SetInt(v)
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set int64 value to %s", ErrInput, val.Kind())
		}

	case Uint64:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set uint64 value", ErrInput)
		}
		var v uint64
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode uint64: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Uint64:
			val.SetUint(v)
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set uint64 value to %s", ErrInput, val.Kind())
		}

	case Float32:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set float32 value", ErrInput)
		}
		var v float32
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode float32: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Float32:
			val.SetFloat(float64(v))
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set float32 value to %s", ErrInput, val.Kind())
		}

	case Float64:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set float64 value", ErrInput)
		}
		var v float64
		if err := binary.Read(buf, binary.LittleEndian, &v); err != nil {
			return fmt.Errorf("%w: failed to decode float64: %w", ErrBuffer, err)
		}

		switch val.Kind() {
		case reflect.Float64:
			val.SetFloat(v)
		case reflect.Interface:
			val.Set(reflect.ValueOf(v))
		default:
			return fmt.Errorf("%w: cannot set float64 value to %s", ErrInput, val.Kind())
		}

	case Bool:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set bool value", ErrInput)
		}
		b, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("%w: failed to decode bool: %w", ErrBuffer, err)
		}

		boolVal := b != 0
		switch val.Kind() {
		case reflect.Bool:
			val.SetBool(boolVal)
		case reflect.Interface:
			val.Set(reflect.ValueOf(boolVal))
		default:
			return fmt.Errorf("%w: cannot set bool value to %s", ErrInput, val.Kind())
		}

	case String:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set string value", ErrInput)
		}

		var length uint32
		if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
			return fmt.Errorf("%w: failed to decode string length: %w", ErrBuffer, err)
		}
		if length > uint32(buf.Len()) {
			return fmt.Errorf("string length %d exceeds buffer size %d", length, buf.Len())
		}

		if length == 0 {
			val.SetString("")
			return nil
		}

		data := buf.Next(int(length))
		if len(data) != int(length) {
			return io.ErrUnexpectedEOF
		}
		strVal := string(data)

		switch val.Kind() {
		case reflect.String:
			val.SetString(strVal)
		case reflect.Interface:
			val.Set(reflect.ValueOf(strVal))
		default:
			return fmt.Errorf("%w: cannot set string value to %s", ErrInput, val.Kind())
		}

	case Bytes:
		if !val.CanSet() {
			return fmt.Errorf("%w: cannot set bytes value", ErrInput)
		}

		var length uint32
		if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
			return fmt.Errorf("%w: failed to decode bytes length: %w", ErrBuffer, err)
		}
		if length > uint32(buf.Len()) {
			return fmt.Errorf("bytes length %d exceeds buffer size %d", length, buf.Len())
		}

		data := make([]byte, length)
		n, err := buf.Read(data)
		if err != nil {
			return fmt.Errorf("%w: failed to decode bytes data: %w", ErrBuffer, err)
		}
		if n != int(length) {
			return fmt.Errorf("%w: expected to read %d bytes, got %d", ErrBuffer, length, n)
		}

		switch val.Kind() {
		case reflect.String:
			val.SetBytes(data)
		case reflect.Interface:
			val.Set(reflect.ValueOf(data))
		default:
			return fmt.Errorf("%w: cannot set bytes value to %s", ErrInput, val.Kind())
		}

	default:
		return fmt.Errorf("%w: unknown SimpleType: %d", ErrModel, t)
	}

	return nil
}
