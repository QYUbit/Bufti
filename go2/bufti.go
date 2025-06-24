package bufti2

import (
	"bytes"
	"errors"
	"sync"
)

const ProtocolVersion uint32 = 1

var (
	// Indicates an incompatible version of a buffer.
	ErrVersion = errors.New("incompatible buffer version")

	// Indicates an invalid encode input value or decode testination value.
	ErrInput = errors.New("unexpected input")

	// Indicates an invalid buffer.
	ErrBuffer = errors.New("invalid buffer")

	// Indicates an error inside of a model.
	ErrModel = errors.New("invalid model")
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

type ModelField struct {
	index      byte
	label      string
	fieldType  BuftiType
	isRequired bool
}

func Field(index byte, label string, fieldType BuftiType) ModelField {
	return ModelField{
		index:     index,
		label:     label,
		fieldType: fieldType,
	}
}

func RequiredField(index byte, label string, fieldType BuftiType) ModelField {
	return ModelField{
		index:      index,
		label:      label,
		fieldType:  fieldType,
		isRequired: true,
	}
}

type Model struct {
	name   string
	schema map[byte]ModelField
	labels map[string]byte
}

// TODO Better model initializer function
func NewModel(name string, fields ...ModelField) *Model {
	m := &Model{
		name:   name,
		schema: make(map[byte]ModelField),
		labels: make(map[string]byte),
	}

	for _, f := range fields {
		m.labels[f.label] = f.index
		m.schema[f.index] = f
	}
	return m
}
