package bufti2

import "testing"

func Test1(t *testing.T) {
	userModel := NewModel("user",
		NewField(0, "id", Int64Type),
		NewField(1, "name", StringType),
		NewField(2, "postIDs", NewList(Int64Type)),
		NewField(3, "avgViews", Float64Type),
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

	t.Log(b, len(b))
}
