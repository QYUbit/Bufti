package bufti2

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 512))
	},
}

type Field struct {
	index     byte
	label     string
	fieldType BuftiType
}

func NewField(index byte, label string, fieldType BuftiType) Field {
	return Field{
		index:     index,
		label:     label,
		fieldType: fieldType,
	}
}

type Model struct {
	name   string
	schema map[byte]Field
	labels map[string]byte
}

func NewModel(name string, fields ...Field) *Model {
	m := &Model{
		name:   name,
		schema: make(map[byte]Field),
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
		NewField(0, "id", Int64Type),
		NewField(1, "name", StringType),
		NewField(2, "postIDs", NewList(Int64Type)),
		NewField(3, "avgViews", Float64Type),
	)

	commentModel := NewModel("comment",
		NewField(0, "id", Int64Type),
		NewField(1, "authorID", Int64Type),
		NewField(1, "postID", Int64Type),
	)

	postModel := NewModel("post",
		NewField(0, "id", Int64Type),
		NewField(1, "title", StringType),
		NewField(2, "tags", NewList(StringType)),
		NewField(3, "views", Int32Type),
		NewField(4, "author", NewReference(userModel)),
		NewField(5, "comments", NewList(NewReference(commentModel))),
	)

	_ = postModel
}
