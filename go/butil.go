package butil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"
)

// ProtocolVersion defines the current version of the binary protocol.
// Buffers encoded with different versions cannot be decoded.
const ProtocolVersion uint32 = 1

var (
	// ErrVersion indicates an incompatible version of a buffer.
	// This occurs when trying to decode data encoded with a different protocol version.
	ErrVersion = errors.New("incompatible buffer version")

	// ErrInput indicates an invalid encode input value or decode destination value.
	// This includes nil inputs, wrong types, or incompatible destination types.
	ErrInput = errors.New("unexpected input")

	// ErrBuffer indicates an invalid or corrupted buffer.
	// This occurs when the binary data cannot be properly read or parsed.
	ErrBuffer = errors.New("invalid buffer")

	// ErrModel indicates an error with the model schema.
	// This includes references to non-existent fields or schema inconsistencies.
	ErrModel = errors.New("invalid model")
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

// ModelField represents a single field in a model schema.
type ModelField struct {
	index      byte
	label      string
	fieldType  BuftiType
	isRequired *bool
}

// Field creates a new model field with the given index, label, and type.
// The field will be required by default unless the model is configured otherwise
func Field(index byte, label string, fieldType BuftiType) ModelField {
	return ModelField{
		index:     index,
		label:     label,
		fieldType: fieldType,
	}
}

// RequiredField creates a new required model field.
// Required fields must be present in the input data during encoding.
func RequiredField(index byte, label string, fieldType BuftiType) ModelField {
	trueValue := true
	return ModelField{
		index:      index,
		label:      label,
		fieldType:  fieldType,
		isRequired: &trueValue,
	}
}

// OptionalField creates a new optional model field.
// Optional fields can be omitted from the input data during encoding.
func OptionalField(index byte, label string, fieldType BuftiType) ModelField {
	falseValue := false
	return ModelField{
		index:      index,
		label:      label,
		fieldType:  fieldType,
		isRequired: &falseValue,
	}
}

// Model represents a schema for binary serialization and deserialization.
type Model struct {
	name       string
	schema     map[byte]ModelField
	labels     map[string]byte
	fieldCache map[reflect.Type]map[string]reflect.Value
	mu         sync.RWMutex
}

// NewModel creates a new model with the given fields.
// All fields are required by default unless explicitly marked as optional.
// Returns ErrModel if the model is misconfigured (such as duplicate indices).
func NewModel(fields ...ModelField) (*Model, error) {
	m := &Model{
		name:   "unnamed_model",
		schema: make(map[byte]ModelField),
		labels: make(map[string]byte),
	}

	for _, f := range fields {
		if f.isRequired == nil {
			trueValue := true
			f.isRequired = &trueValue
		}
		m.labels[f.label] = f.index
		m.schema[f.index] = f
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// ModelOptions configures how a model is created.
type ModelOptions struct {
	// Name is a human-readable identifier for the model, used in error messages.
	Name string
	// RequiredByDefault determines whether fields without explicit required/optional
	RequiredByDefault bool
}

// NewModelWithOptions creates a new model with the given options and fields.
// This allows customization of the model's behavior.
// Returns ErrModel if the model is misconfigured (such as duplicate indices).
func NewModelWithOptions(options *ModelOptions, fields ...ModelField) (*Model, error) {
	m := &Model{
		name:   options.Name,
		schema: make(map[byte]ModelField),
		labels: make(map[string]byte),
	}

	for _, f := range fields {
		if f.isRequired == nil {
			if options.RequiredByDefault {
				trueValue := true
				f.isRequired = &trueValue
			} else {
				falseValue := false
				f.isRequired = &falseValue
			}
		}
		m.labels[f.label] = f.index
		m.schema[f.index] = f
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Model) Validate() error {
	var labels []string
	for label := range m.labels {
		if slices.Contains(labels, label) {
			return fmt.Errorf("%w: duplicate label %s", ErrModel, label)
		}
		labels = append(labels, label)
	}

	var indices []byte
	for index := range m.schema {
		if slices.Contains(indices, index) {
			return fmt.Errorf("%w: duplicate index %d", ErrModel, index)
		}
		indices = append(indices, index)
	}

	return nil
}

// Used for unit testing
func newModel(fields ...ModelField) *Model {
	model, _ := NewModel(fields...)
	return model
}

func newModelWithOptions(options *ModelOptions, fields ...ModelField) *Model {
	model, _ := NewModelWithOptions(options, fields...)
	return model
}
