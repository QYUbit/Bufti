package bufti2

import (
	"testing"
)

var userModel = NewModel(
	Field(0, "id", Int64),
	Field(1, "name", String),
	Field(2, "postIDs", List(Int64)),
	Field(3, "avgViews", Float64),
)

func TestStruct(t *testing.T) {
	type User struct {
		Id       int     `bufti:"id"`
		Name     string  `bufti:"name"`
		Posts    []int64 `bufti:"postIDs"`
		AvgViews float64 `bufti:"avgViews"`
	}

	userStruct1 := &User{
		Id:       238476294,
		Name:     "alice",
		Posts:    []int64{6884623391, 2382322946, 3674892946},
		AvgViews: 23.8,
	}

	b, err := userModel.Encode(userStruct1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)

	var userStruct2 User
	if err := userModel.Decode(b, &userStruct2); err != nil {
		t.Fatal(err)
	}

	t.Log(userStruct2)
}

func TestMap(t *testing.T) {
	userMap1 := map[string]any{
		"id":       2384762946,
		"name":     "alice",
		"postIDs":  []int64{6884623391, 2382322946, 3674892946},
		"avgViews": 23.8,
	}

	b, err := userModel.Encode(userMap1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)

	userMap2 := make(map[string]any)
	if err := userModel.Decode(b, &userMap2); err != nil {
		t.Fatal(err)
	}

	t.Log(userMap2)
}

func TestMapType(t *testing.T) {
	type MapStruct struct {
		M map[string]int `bufti:"map"`
	}

	mapModel := NewModel(
		Field(0, "map", Map(String, Int64)),
	)

	mapStruct := &MapStruct{
		M: map[string]int{
			"a": 1,
			"b": 2,
			"c": 4,
		},
	}

	b, err := mapModel.Encode(mapStruct)
	if err != nil {
		t.Fatal(b)
	}

	t.Log(b)

	var dest MapStruct
	if err := mapModel.Decode(b, &dest); err != nil {
		t.Fatal(err)
	}

	t.Log(dest)
}

func TestMapMap(t *testing.T) {
	mapModel := NewModel(
		Field(0, "map", Map(String, Int64)),
	)

	mapmap := map[string]any{
		"map": map[string]int{
			"a": 1,
			"b": 2,
			"c": 4,
		},
	}

	b, err := mapModel.Encode(mapmap)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)

	dest := make(map[string]any)
	if err := mapModel.Decode(b, &dest); err != nil {
		t.Fatal(err)
	}

	t.Log(dest)
}

func TestReferenceType(t *testing.T) { // ! Doesn't work
	type OtherStruct struct {
		Text string `bufti:"text"`
	}

	type ReferenceStruct struct {
		Reference *OtherStruct `bufti:"reference"`
	}

	otherModel := NewModel(
		Field(0, "text", String),
	)

	referenceModel := NewModel(
		Field(0, "reference", Reference(otherModel)),
	)

	testStruct1 := &ReferenceStruct{
		Reference: &OtherStruct{
			Text: "testestest",
		},
	}

	b, err := referenceModel.Encode(testStruct1)
	if err != nil {
		t.Fatal(b)
	}

	t.Log(b)

	var dest ReferenceStruct
	if err := referenceModel.Decode(b, &dest); err != nil {
		t.Fatal(err)
	}

	t.Log(dest)
}
