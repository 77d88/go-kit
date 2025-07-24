package xcore

import (
	"fmt"
	"testing"
)

func TestIsZeroValue(t *testing.T) {
	var s *string
	fmt.Println(IsZero(s)) // 输出: true

	s = new(string)
	fmt.Println(IsZero(s)) // 输出: true

	x := "not empty"
	s = &x
	fmt.Println(IsZero(s)) // 输出: false
	fmt.Println(IsZero(x)) // 输出: false

	var num *int
	fmt.Println(IsZero(num)) // 输出: true

	num = new(int)
	fmt.Println(IsZero(num)) // 输出: true

	*num = 10
	fmt.Println(IsZero(num)) // 输出: false
	type A struct {
		A string
	}

	var a A
	var ax *A
	fmt.Println(IsZero(A{}))
	fmt.Println(IsZero(new(A)))
	fmt.Println(IsZero(a))
	fmt.Println(IsZero(ax))
	fmt.Println(IsZero(A{A: "1"}))
}

func TestTR(t *testing.T) {
	ternaryFunc := TernaryFunc(false, func() string {
		return "true"
	}, func() string {
		return "false"
	})
	fmt.Println(ternaryFunc)
}

func TestName(t *testing.T) {
	a := struct {
		A string
	}{
		A: "123",
	}
	s, err := Copy(a)
	fmt.Printf("%v===%v\n", s, err)
}
