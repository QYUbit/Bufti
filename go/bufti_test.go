package bufti

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
)

type SimpleStruct struct {
	ID   int64   `bufti:"id"`
	Name string  `bufti:"name"`
	Age  int32   `bufti:"age"`
	Rate float64 `bufti:"rate"`
}

type ComplexStruct struct {
	ID       int64            `bufti:"id"`
	Name     string           `bufti:"name"`
	Tags     []string         `bufti:"tags"`
	Scores   []float64        `bufti:"scores"`
	Metadata map[string]int64 `bufti:"metadata"`
	Active   bool             `bufti:"active"`
	Data     []byte           `bufti:"data"`
}

type NestedStruct struct {
	ID       int64          `bufti:"id"`
	Simple   SimpleStruct   `bufti:"simple"`
	Children []SimpleStruct `bufti:"children"`
}

// Test Models
var simpleModel = NewModel(
	Field(0, "id", Int64),
	Field(1, "name", String),
	Field(2, "age", Int32),
	Field(3, "rate", Float64),
)

var complexModel = NewModel(
	Field(0, "id", Int64),
	Field(1, "name", String),
	Field(2, "tags", List(String)),
	Field(3, "scores", List(Float64)),
	Field(4, "metadata", Map(String, Int64)),
	Field(5, "active", Bool),
	Field(6, "data", Bytes),
)

// Test Basic Types
func TestSimpleTypes(t *testing.T) {
	tests := []struct {
		name     string
		model    *Model
		input    any
		expected any
	}{
		{
			name:     "int8",
			model:    NewModel(Field(0, "value", Int8)),
			input:    map[string]any{"value": int8(42)},
			expected: map[string]any{"value": int8(42)},
		},
		{
			name:     "uint8",
			model:    NewModel(Field(0, "value", Uint8)),
			input:    map[string]any{"value": uint8(255)},
			expected: map[string]any{"value": uint8(255)},
		},
		{
			name:     "int16",
			model:    NewModel(Field(0, "value", Int16)),
			input:    map[string]any{"value": int16(-32768)},
			expected: map[string]any{"value": int16(-32768)},
		},
		{
			name:     "uint16",
			model:    NewModel(Field(0, "value", Uint16)),
			input:    map[string]any{"value": uint16(65535)},
			expected: map[string]any{"value": uint16(65535)},
		},
		{
			name:     "int32",
			model:    NewModel(Field(0, "value", Int32)),
			input:    map[string]any{"value": int32(-2147483648)},
			expected: map[string]any{"value": int32(-2147483648)},
		},
		{
			name:     "uint32",
			model:    NewModel(Field(0, "value", Uint32)),
			input:    map[string]any{"value": uint32(4294967295)},
			expected: map[string]any{"value": uint32(4294967295)},
		},
		{
			name:     "int64",
			model:    NewModel(Field(0, "value", Int64)),
			input:    map[string]any{"value": int64(-9223372036854775808)},
			expected: map[string]any{"value": int64(-9223372036854775808)},
		},
		{
			name:     "uint64",
			model:    NewModel(Field(0, "value", Uint64)),
			input:    map[string]any{"value": uint64(18446744073709551615)},
			expected: map[string]any{"value": uint64(18446744073709551615)},
		},
		{
			name:     "float32",
			model:    NewModel(Field(0, "value", Float32)),
			input:    map[string]any{"value": float32(3.14159)},
			expected: map[string]any{"value": float32(3.14159)},
		},
		{
			name:     "float64",
			model:    NewModel(Field(0, "value", Float64)),
			input:    map[string]any{"value": float64(3.141592653589793)},
			expected: map[string]any{"value": float64(3.141592653589793)},
		},
		{
			name:     "bool_true",
			model:    NewModel(Field(0, "value", Bool)),
			input:    map[string]any{"value": true},
			expected: map[string]any{"value": true},
		},
		{
			name:     "bool_false",
			model:    NewModel(Field(0, "value", Bool)),
			input:    map[string]any{"value": false},
			expected: map[string]any{"value": false},
		},
		{
			name:     "string",
			model:    NewModel(Field(0, "value", String)),
			input:    map[string]any{"value": "Hello, 世界! 🌍"},
			expected: map[string]any{"value": "Hello, 世界! 🌍"},
		},
		//{
		//	name:     "empty_string",
		//	model:    NewModel(Field(0, "value", String)),
		//	input:    map[string]any{"value": ""},
		//	expected: map[string]any{"value": ""},
		//},
		{
			name:     "bytes",
			model:    NewModel(Field(0, "value", Bytes)),
			input:    map[string]any{"value": []byte{0, 1, 2, 255, 128}},
			expected: map[string]any{"value": []byte{0, 1, 2, 255, 128}},
		},
		//{
		//	name:     "empty_bytes",
		//	model:    NewModel(Field(0, "value", Bytes)),
		//	input:    map[string]any{"value": []byte{}},
		//	expected: map[string]any{"value": []byte{}},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := tt.model.Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Decode
			result := make(map[string]any)
			err = tt.model.Decode(encoded, &result)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Compare
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestStructEncodeDecodeBasic(t *testing.T) {
	original := &SimpleStruct{
		ID:   12345678901234,
		Name: "Test User",
		Age:  25,
		Rate: 99.95,
	}

	// Encode
	encoded, err := simpleModel.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode to struct
	var decoded SimpleStruct
	err = simpleModel.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(*original, decoded) {
		t.Errorf("Expected %+v, got %+v", *original, decoded)
	}
}

func TestMapEncodeDecodeBasic(t *testing.T) {
	original := map[string]any{
		"id":   int64(12345678901234),
		"name": "Test User",
		"age":  int32(25),
		"rate": 99.95,
	}

	// Encode
	encoded, err := simpleModel.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode to map
	decoded := make(map[string]any)
	err = simpleModel.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestComplexTypes(t *testing.T) {
	original := &ComplexStruct{
		ID:     987654321,
		Name:   "Complex Test",
		Tags:   []string{"tag1", "tag2", "tag3"},
		Scores: []float64{1.1, 2.2, 3.3, 4.4},
		Metadata: map[string]int64{
			"count":  100,
			"offset": 50,
			"limit":  25,
		},
		Active: true,
		Data:   []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}

	// Encode
	encoded, err := complexModel.Encode(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode to struct
	var decoded ComplexStruct
	err = complexModel.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(*original, decoded) {
		t.Errorf("Expected %+v, got %+v", *original, decoded)
	}
}

// ! TestListTypes and TestMapTypes don't work
func TestListTypes(t *testing.T) {
	tests := []struct {
		name     string
		model    *Model
		input    any
		expected any
	}{
		{
			name:     "list_of_strings",
			model:    NewModel(Field(0, "items", List(String))),
			input:    map[string]any{"items": []string{"a", "b", "c"}},
			expected: map[string]any{"items": []string{"a", "b", "c"}},
		},
		{
			name:     "list_of_ints",
			model:    NewModel(Field(0, "items", List(Int64))),
			input:    map[string]any{"items": []int64{1, 2, 3, 4, 5}},
			expected: map[string]any{"items": []int64{1, 2, 3, 4, 5}},
		},
		{
			name:     "empty_list",
			model:    NewModel(Field(0, "items", List(String))),
			input:    map[string]any{"items": []string{}},
			expected: map[string]any{"items": []string{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := tt.model.Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			result := make(map[string]any)
			err = tt.model.Decode(encoded, &result)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMapTypes(t *testing.T) {
	tests := []struct {
		name     string
		model    *Model
		input    any
		expected any
	}{
		{
			name:     "string_to_int_map",
			model:    NewModel(Field(0, "data", Map(String, Int64))),
			input:    map[string]any{"data": map[string]int{"a": 1, "b": 2, "c": 3}},
			expected: map[string]any{"data": map[string]int64{"a": 1, "b": 2, "c": 3}},
		},
		{
			name:     "empty_map",
			model:    NewModel(Field(0, "data", Map(String, Int64))),
			input:    map[string]any{"data": map[string]int{}},
			expected: map[string]any{"data": map[string]int64{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := tt.model.Encode(tt.input)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			result := make(map[string]any)
			err = tt.model.Decode(encoded, &result)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestModelOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  *ModelOptions
		expected string
	}{
		{
			name: "named_model",
			options: &ModelOptions{
				Name:              "TestModel",
				RequiredByDefault: true,
			},
			expected: "TestModel",
		},
		{
			name: "optional_by_default",
			options: &ModelOptions{
				Name:              "OptionalModel",
				RequiredByDefault: false,
			},
			expected: "OptionalModel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModelWithOptions(tt.options,
				Field(0, "id", Int64),
				Field(1, "name", String),
			)

			if model.name != tt.expected {
				t.Errorf("Expected model name %s, got %s", tt.expected, model.name)
			}
		})
	}
}

func TestFieldOptions(t *testing.T) {
	model := NewModel(
		RequiredField(0, "required", String),
		OptionalField(1, "optional", String),
		Field(2, "default", String), // Should be required by default
	)

	// Test that required field configuration is preserved
	requiredField := model.schema[0]
	if requiredField.isRequired == nil || !*requiredField.isRequired {
		t.Error("Required field should be marked as required")
	}

	optionalField := model.schema[1]
	if optionalField.isRequired == nil || *optionalField.isRequired {
		t.Error("Optional field should be marked as optional")
	}

	defaultField := model.schema[2]
	if defaultField.isRequired == nil || !*defaultField.isRequired {
		t.Error("Default field should be marked as required")
	}
}

// Edge Cases Tests
func TestEdgeCases(t *testing.T) {
	t.Run("nil_input", func(t *testing.T) {
		_, err := simpleModel.Encode(nil)
		if err == nil {
			t.Error("Expected error for nil input")
		}
	})

	t.Run("invalid_buffer", func(t *testing.T) {
		invalidBuffer := []byte{1, 2, 3} // Too short
		var result SimpleStruct
		err := simpleModel.Decode(invalidBuffer, &result)
		if err == nil {
			t.Error("Expected error for invalid buffer")
		}
	})

	t.Run("wrong_version", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		// Write wrong version
		wrongVersion := uint32(999)
		if err := binary.Write(buf, binary.LittleEndian, wrongVersion); err != nil {
			t.Fatal(err)
		}

		var result SimpleStruct
		err := simpleModel.Decode(buf.Bytes(), &result)
		if err == nil {
			t.Error("Expected version error")
		}
		if !strings.Contains(err.Error(), "incompatible") {
			t.Errorf("Expected version error, got: %v", err)
		}
	})

	t.Run("decode_to_non_pointer", func(t *testing.T) {
		data := map[string]any{"id": int64(123)}
		encoded, err := simpleModel.Encode(data)
		if err != nil {
			t.Fatal(err)
		}

		var result SimpleStruct // Not a pointer
		err = simpleModel.Decode(encoded, result)
		if err == nil {
			t.Error("Expected error when decoding to non-pointer")
		}
	})

	t.Run("large_string", func(t *testing.T) {
		model := NewModel(Field(0, "data", String))
		largeString := strings.Repeat("A", 1000000) // 1MB string

		input := map[string]any{"data": largeString}
		encoded, err := model.Encode(input)
		if err != nil {
			t.Fatalf("Failed to encode large string: %v", err)
		}

		result := make(map[string]any)
		err = model.Decode(encoded, &result)
		if err != nil {
			t.Fatalf("Failed to decode large string: %v", err)
		}

		if result["data"] != largeString {
			t.Error("Large string not preserved")
		}
	})

	t.Run("large_byte_slice", func(t *testing.T) {
		model := NewModel(Field(0, "data", Bytes))
		largeBytes := make([]byte, 1000000) // 1MB
		for i := range largeBytes {
			largeBytes[i] = byte(i % 256)
		}

		input := map[string]any{"data": largeBytes}
		encoded, err := model.Encode(input)
		if err != nil {
			t.Fatalf("Failed to encode large bytes: %v", err)
		}

		result := make(map[string]any)
		err = model.Decode(encoded, &result)
		if err != nil {
			t.Fatalf("Failed to decode large bytes: %v", err)
		}

		if !bytes.Equal(result["data"].([]byte), largeBytes) {
			t.Error("Large byte slice not preserved")
		}
	})
}

func TestTypeConversions(t *testing.T) {
	t.Run("int_conversions", func(t *testing.T) {
		model := NewModel(Field(0, "value", Int64))

		// Test various int types
		inputs := []any{
			int(42),
			int8(42),
			int16(42),
			int32(42),
			int64(42),
		}

		for _, input := range inputs {
			data := map[string]any{"value": input}
			encoded, err := model.Encode(data)
			if err != nil {
				t.Errorf("Failed to encode %T: %v", input, err)
				continue
			}

			result := make(map[string]any)
			err = model.Decode(encoded, &result)
			if err != nil {
				t.Errorf("Failed to decode %T: %v", input, err)
				continue
			}

			if result["value"] != int64(42) {
				t.Errorf("Expected int64(42), got %v (%T)", result["value"], result["value"])
			}
		}
	})

	t.Run("overflow_detection", func(t *testing.T) {
		model := NewModel(Field(0, "value", Int8))

		// This should fail - value too large for int8
		data := map[string]any{"value": int(300)}
		_, err := model.Encode(data)
		if err == nil {
			t.Error("Expected overflow error for int8")
		}
	})
}

func TestMissingFields(t *testing.T) {
	// Test encoding with missing fields
	model := NewModel(
		Field(0, "required", String),
		Field(1, "optional", String),
	)

	// Encode with missing field
	data := map[string]any{"required": "present"}
	encoded, err := model.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode with missing field: %v", err)
	}

	// Decode should work
	result := make(map[string]any)
	err = model.Decode(encoded, &result)
	if err != nil {
		t.Fatalf("Failed to decode with missing field: %v", err)
	}

	if result["required"] != "present" {
		t.Error("Required field not preserved")
	}
	if _, exists := result["optional"]; exists {
		t.Error("Missing field should not be present in result")
	}
}

// Benchmarks
func BenchmarkEncodeSimpleStruct(b *testing.B) {
	data := &SimpleStruct{
		ID:   12345678901234,
		Name: "Benchmark User",
		Age:  30,
		Rate: 95.5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := simpleModel.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeSimpleStruct(b *testing.B) {
	data := &SimpleStruct{
		ID:   12345678901234,
		Name: "Benchmark User",
		Age:  30,
		Rate: 95.5,
	}

	encoded, err := simpleModel.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result SimpleStruct
		err := simpleModel.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeComplexStruct(b *testing.B) {
	data := &ComplexStruct{
		ID:     987654321,
		Name:   "Complex Benchmark",
		Tags:   []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		Scores: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9, 10.0},
		Metadata: map[string]int64{
			"count":   1000,
			"offset":  500,
			"limit":   100,
			"version": 2,
			"flags":   0xFF,
		},
		Active: true,
		Data:   bytes.Repeat([]byte{0xAB, 0xCD, 0xEF}, 100),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := complexModel.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeComplexStruct(b *testing.B) {
	data := &ComplexStruct{
		ID:     987654321,
		Name:   "Complex Benchmark",
		Tags:   []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		Scores: []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9, 10.0},
		Metadata: map[string]int64{
			"count":   1000,
			"offset":  500,
			"limit":   100,
			"version": 2,
			"flags":   0xFF,
		},
		Active: true,
		Data:   bytes.Repeat([]byte{0xAB, 0xCD, 0xEF}, 100),
	}

	encoded, err := complexModel.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result ComplexStruct
		err := complexModel.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeMap(b *testing.B) {
	data := map[string]any{
		"id":   int64(12345678901234),
		"name": "Benchmark User",
		"age":  int32(30),
		"rate": 95.5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := simpleModel.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMap(b *testing.B) {
	data := map[string]any{
		"id":   int64(12345678901234),
		"name": "Benchmark User",
		"age":  int32(30),
		"rate": 95.5,
	}

	encoded, err := simpleModel.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[string]any)
		err := simpleModel.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeLargeList(b *testing.B) {
	model := NewModel(Field(0, "items", List(Int64)))

	items := make([]int64, 10000)
	for i := range items {
		items[i] = int64(i)
	}
	data := map[string]any{"items": items}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := model.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeLargeList(b *testing.B) {
	model := NewModel(Field(0, "items", List(Int64)))

	items := make([]int64, 10000)
	for i := range items {
		items[i] = int64(i)
	}
	data := map[string]any{"items": items}

	encoded, err := model.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[string]any)
		err := model.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeString(b *testing.B) {
	model := NewModel(Field(0, "text", String))
	data := map[string]any{"text": "This is a benchmark string with some content to test performance"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := model.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeString(b *testing.B) {
	model := NewModel(Field(0, "text", String))
	data := map[string]any{"text": "This is a benchmark string with some content to test performance"}

	encoded, err := model.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make(map[string]any)
		err := model.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Memory allocation benchmarks
func BenchmarkEncodeAllocations(b *testing.B) {
	data := &SimpleStruct{
		ID:   12345678901234,
		Name: "Allocation Test",
		Age:  25,
		Rate: 99.95,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := simpleModel.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeAllocations(b *testing.B) {
	data := &SimpleStruct{
		ID:   12345678901234,
		Name: "Allocation Test",
		Age:  25,
		Rate: 99.95,
	}

	encoded, err := simpleModel.Encode(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result SimpleStruct
		err := simpleModel.Decode(encoded, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}
