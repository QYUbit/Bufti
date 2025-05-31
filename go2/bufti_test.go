package bufti2

import (
	"testing"
)

var userModel = NewModel("user",
	Field(0, "id", Int64),
	Field(1, "name", String),
	Field(2, "postIDs", List(Int64)),
	Field(3, "avgViews", Float64),
)

type User struct {
	Id       int     `bufti:"id"`
	Name     string  `bufti:"name"`
	Posts    []int64 `bufti:"postIDs"`
	AvgViews float64 `bufti:"avgViews"`
}

func TestStruct(t *testing.T) {
	userStruct1 := &User{
		Id:       2384762946,
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
