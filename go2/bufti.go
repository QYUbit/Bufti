package bufti2

import (
	"bytes"
	"sync"
)

const ProtocolVersion uint32 = 2

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

type ModelField struct {
	index     byte
	label     string
	fieldType BuftiType
}

func Field(index byte, label string, fieldType BuftiType) ModelField {
	return ModelField{
		index:     index,
		label:     label,
		fieldType: fieldType,
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

func M() {
	userModel := NewModel("user",
		Field(0, "id", Int64),
		Field(1, "name", String),
		Field(2, "postIDs", List(Int64)),
		Field(3, "avgViews", Float64),
	)

	commentModel := NewModel("comment",
		Field(0, "id", Int64),
		Field(1, "authorID", Int64),
		Field(1, "postID", Int64),
	)

	postModel := NewModel("post",
		Field(0, "id", Int64),
		Field(1, "title", String),
		Field(2, "tags", List(String)),
		Field(3, "views", Int32),
		Field(4, "author", Reference(userModel)),
		Field(5, "comments", List(Reference(commentModel))),
	)

	_ = postModel
}
