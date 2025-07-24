package xcore

import (
	"reflect"
)

// IsZero 检查值是否等于其类型的零值
func IsZero[T any](val T) bool {
	switch v := any(val).(type) {
	// 基础类型处理
	case bool:
		return !v
	case int, int8, int16, int32, int64:
		return any(v) == 0
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return any(v) == 0
	case float32, float64:
		return any(v) == 0
	case string:
		return v == ""
	case []byte:
		return len(v) == 0

	// 指针类型特殊处理
	case *bool:
		return v == nil || !*v
	case *int:
		return v == nil || *v == 0
	case *int8:
		return v == nil || *v == 0
	case *int16:
		return v == nil || *v == 0
	case *int32:
		return v == nil || *v == 0
	case *int64:
		return v == nil || *v == 0
	case *uint:
		return v == nil || *v == 0
	case *uint8:
		return v == nil || *v == 0
	case *uint16:
		return v == nil || *v == 0
	case *uint32:
		return v == nil || *v == 0
	case *uint64:
		return v == nil || *v == 0
	case *uintptr:
		return v == nil || *v == 0
	case *float32:
		return v == nil || *v == 0
	case *float64:
		return v == nil || *v == 0
	case *string:
		return v == nil || *v == ""
	case *[]byte:
		return v == nil || len(*v) == 0

	// 其他类型使用反射处理
	default:
		rv := reflect.ValueOf(val)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return true
			}
			return reflect.DeepEqual(rv.Elem().Interface(), reflect.Zero(rv.Elem().Type()).Interface())
		}
		return reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface())
	}
}

// IsBasicType 检查给定值是否为基本类型。
func IsBasicType(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool,
		string:
		return true
	default:
		return false
	}
}

// Ternary 返回一个布尔值，如果 condition 为 true，则返回 trueVal，否则返回 falseVal。
func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// TernaryFunc 返回一个布尔值，如果 condition 为 true，则返回 trueVal()，否则返回 falseVal()。
func TernaryFunc[R any, T func() R](condition bool, trueVal, falseVal T) R {
	if condition {
		return trueVal()
	}
	return falseVal()
}

func NewBy(src interface{}) interface{} {
	value := reflect.ValueOf(src)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return reflect.New(value.Type()).Interface()
}

// FirstOrDefault 返回切片中的第一个元素，如果切片为空则返回默认值
// 参数:
//
//	defaultValue: 当切片为空时返回的默认值
//	values: 要检查的切片，可以是零个或多个元素
//
// 返回值:
//
//	切片中的第一个元素，或默认值(当切片为空时)
//
// 使用示例:
//
//	result := FirstOrDefault(0, 1, 2, 3) // 返回1
//	result := FirstOrDefault("default")  // 返回"default"
func FirstOrDefault[T any](defaultValue T, values ...T) T {
	if len(values) == 0 {
		return defaultValue
	}
	return values[0]
}

// V2p 转换并返回一个类型的值指针。
// T: 任何类型。
// 返回类型 *T: 类型T的指针。
// V2p 赋值对象并返回复制体的指针
func V2p[T any](t T) *T {
	return &t
}

// P2v 转换并返回一个类型指针的值。
// T: 任何类型。
// t: 类型T的指针。
// 返回类型 T: 类型T的值。
// P2v 返回指针对应值
func P2v[T any](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}
