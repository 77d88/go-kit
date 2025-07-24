package xstr

import (
	"fmt"
	"testing"
)

func TestToInt(t *testing.T) {
	s := "1"

	x, _ := ToInt[int64](s)
	x1, _ := ToInt[int](s)
	fmt.Printf("%d", x+int64(x1))
}

func TestSplitToInt(t *testing.T) {
	s := "1,2,3,4"
	toInt := SplitToInt[int64](s, ",")
	fmt.Printf("%v", toInt)
}

func TestTelMask(t *testing.T) {

	fmt.Printf("%s xxx", TelNumberMask("028-1", "*"))
}
