package xrandom

import (
	"math"
	"math/rand"
)

// RandInt 生成一个随机数 [min, max) 包含 min 不包含 max
func RandInt(min, max int) int {
	if min == max {
		return min
	}

	if max < min {
		min, max = max, min
	}

	if min == 0 && max == math.MaxInt {
		return rand.Int()
	}

	return rand.Intn(max-min) + min
}

// RandBool 随机bool
func RandBool() bool {
	return rand.Intn(2) == 0
}

// RandSlice 从slice中随机返回一个元素
func RandSlice[T any](slice []T) T {
	return slice[RandInt(0, len(slice))]
}

func RandFloat64(min, max float64) float64 {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}

	result := min + rand.Float64()*(max-min)
	if result >= max {
		// 防止因浮点精度越界
		result = max - 1e-12
	}
	return result
}
