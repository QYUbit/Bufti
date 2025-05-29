package bufti2

import (
	"testing"
)

func Test1(t *testing.T) {
	userModel := NewModel("user",
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

	user := &User{
		Id:       2384762946,
		Name:     "alice",
		Posts:    []int64{6884623391, 2382322946, 3674892946},
		AvgViews: 23.8,
	}

	b, err := userModel.Encode(user)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(b)

	var user2 User

	decoder := NewDecoder(b)
	defer decoder.Close()

	if err := decoder.Decode(&user2); err != nil {
		t.Fatal(err)
	}

	t.Log(user2)
}
