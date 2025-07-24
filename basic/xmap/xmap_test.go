package xmap

import (
	"maps"
	"slices"
	"testing"
)

func Test_Keys(t *testing.T) {

	m := map[string]int{"a": 1, "b": 2}
	keys := Keys(m)
	t.Log(keys)

	t.Log(slices.Collect(maps.Keys(m)))
}
