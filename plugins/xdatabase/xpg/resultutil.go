package xpg

import (
	"fmt"
	"reflect"
	"strconv"
)

func mapFirstCovert(m []map[string]any, dest interface{}) (interface{}, bool) {
	var v any
	if len(m) > 0 {
		for _, val := range m[0] {
			v = val
			break
		}
	}
	if v == nil {
		return nil, false
	}

	// 获取目标类型的反射类型
	destType := reflect.TypeOf(dest).Elem()

	// 尝试直接转换
	if reflect.TypeOf(v).AssignableTo(destType) {
		return v, true
	}

	// 处理数值类型转换
	switch destType.Kind() {
	case reflect.Int:
		switch val := v.(type) {
		case int64:
			return int(val), true
		case float64:
			return int(val), true
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				return i, true
			}
		}
	case reflect.Int64:
		switch val := v.(type) {
		case int:
			return int64(val), true
		case float64:
			return int64(val), true
		case string:
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				return i, true
			}
		}
	case reflect.Float64:
		switch val := v.(type) {
		case int:
			return float64(val), true
		case int64:
			return float64(val), true
		case string:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				return f, true
			}
		}
	case reflect.String:
		return fmt.Sprintf("%v", v), true
	case reflect.Bool:
		switch val := v.(type) {
		case int:
			return val != 0, true
		case int64:
			return val != 0, true
		case string:
			if b, err := strconv.ParseBool(val); err == nil {
				return b, true
			}
		}
	default:
		return nil, false
	}

	return nil, false
}
