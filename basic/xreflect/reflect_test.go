package xreflect

import "testing"

type User struct {
	Name string
	Age  int
}

func TestName(t *testing.T) {

	obj := User{
		Name: "xiaoming",
		Age:  18,
	}
	fields := GetAllFields(&obj)

	for name, field := range fields {
		t.Log(name, field.GetVal())
	}
}

func TestSetFieldValue(t *testing.T) {
	obj := User{
		Name: "xiaoming",
		Age:  18,
	}
	t.Log(obj)
	SetFieldValue(&obj, "Name", "xiaohong2")
	t.Log(obj)
}

func TestGetFieldValue(t *testing.T) {
	obj := User{
		Name: "xiaoming",
		Age:  18,
	}
	t.Log(obj)

	t.Log(GetFieldVal(&obj, "Name2"))
}
