package xarray

import (
	"testing"
)

func TestFirstOrDefault(t *testing.T) {
	t.Log(FirstOrDefault([]int{1, 2, 3}, 99))
	t.Log(FirstOrDefault([]int{}, 99))
	t.Log(FirstOrDefault([]int{0, 2, 3}, 99, true))
	t.Log(FirstOrDefault([]string{""}, "99", true))
	t.Log(FirstOrDefault([]string{"123"}, "99", true))
	t.Log(FirstOrDefault([]string{"123"}, "99", false))
	t.Log(FirstOrDefault([]string{""}, "99", false))

}
