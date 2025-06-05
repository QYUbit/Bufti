package bufti2

import (
	"bytes"
	"errors"
	"sync"
)

const ProtocolVersion uint32 = 1

var (
	ErrInput  = errors.New("unexpected input")
	ErrBuffer = errors.New("invalid buffer")
	ErrModel  = errors.New("invalid model")
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
