package xreflect

import (
	"reflect"
	"testing"
)

func TestGetFieldVal(t *testing.T) {
	type args[T any] struct {
		obj   T
		field string
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want interface{}
	}
	tests := []testCase[User]{
		{
			name: "test1",
			args: args[User]{
				obj: User{
					Name: "test",
					Age:  18,
				},
				field: "Name",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFieldVal(tt.args.obj, tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldVal() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestIsSlice(t *testing.T) {
	t.Log(IsSlice([]int{1, 2, 3}))
}

func Test_toSlice(t *testing.T) {
	t.Log(ToSlice(&[]int{1, 2, 3}))

}

func TestIsPointer(t *testing.T) {
	t.Log(IsPointer(nil))
}
