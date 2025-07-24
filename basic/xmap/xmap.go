package xmap

import (
	"maps"
	"reflect"
	"slices"
)

// Keys 获取 map 中所有的键
func Keys[T comparable, B any](m map[T]B) []T {
	if m == nil {
		return make([]T, 0)
	}
	return slices.Collect(maps.Keys(m))
}

// Values 获取 map 中所有的值
func Values[T comparable, B any](m map[T]B) []B {
	if m == nil {
		return make([]B, 0)
	}
	return slices.Collect(maps.Values(m))
}

// ForEach 遍历 map
func ForEach[T comparable, B any, C any](m map[T]B, iteratee func(key T, value B) C) {
	if m == nil {
		return
	}
	for k, v := range m {
		iteratee(k, v)
	}
}

// IsMap 判断是否 map
func IsMap[T any](obj T) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Map
}
