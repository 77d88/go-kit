package xparse

import (
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"reflect"
	"strconv"
	"time"
)

type Parser = func(any) (any, error)

// 自定义解析器注册表
var parsers = make(map[reflect.Type]Parser)

// RegisterParser 注册自定义类型解析器
func RegisterParser[T any](parser Parser) {
	tType := reflect.TypeOf((*T)(nil)).Elem()
	parsers[tType] = parser
}

// ParseCustom 自定义解析函数模板
func ParseCustom[T any](input any) (T, error) {
	var zero T
	if input == nil {
		return zero, xerror.New("nil input")
	}

	tType := reflect.TypeOf((*T)(nil)).Elem()
	if parser, exists := parsers[tType]; exists {
		result, err := parser(input)
		if err != nil {
			return zero, err
		}
		return result.(T), nil
	}

	return zero, xerror.New("no parser found for type " + tType.Name())
}

// ToNumber 字符串转数字（支持泛型）
func ToNumber[T xtype.Numer](s string) (T, error) {
	if xcore.IsZero(s) {
		return 0, xerror.New("empty string")
	}

	var zero T
	switch any(zero).(type) {
	case int, int8, int16, int32, int64:
		v, err := strconv.ParseInt(s, 10, 64)
		return T(v), err
	case uint, uint8, uint16, uint32, uint64:
		v, err := strconv.ParseUint(s, 10, 64)
		return T(v), err
	case float32, float64:
		v, err := strconv.ParseFloat(s, 64)
		return T(v), err
	default:
		return zero, xerror.New("unsupported number type")
	}
}

// ToString 任意类型转字符串
// ToString 任意类型转字符串（优化版）
func ToString[T any](v T, defaultValue ...T) string {
	if xcore.IsZero(v) {
		if len(defaultValue) > 0 {
			return ToString[T](defaultValue[0])
		}
		return ""
	}

	switch val := any(v).(type) {
	case string:
		return val
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToTime 字符串转时间
func ToTime(s string, layout string) (time.Time, error) {
	if xcore.IsZero(s) {
		return time.Time{}, xerror.New("empty string")
	}
	return time.Parse(layout, s)
}

// TimeToString 时间转字符串
func TimeToString(t time.Time, layout string) string {
	if xcore.IsZero(t) {
		return ""
	}
	return t.Format(layout)
}

// ToJSON 结构体转JSON字符串
func ToJSON[T any](v T) (string, error) {
	if xcore.IsZero(v) {
		return "", xerror.New("empty value")
	}
	b, err := json.Marshal(v)
	return string(b), err
}

// FromJSON JSON字符串转结构体
func FromJSON[T any](s string) (T, error) {
	var zero T
	if xcore.IsZero(s) {
		return zero, xerror.New("empty string")
	}
	var v T
	err := json.Unmarshal([]byte(s), &v)
	return v, err
}

// GetInst 获取指向实例
func GetInst[T any](t T) (reflect.Value, reflect.Kind) {
	value := reflect.ValueOf(t)
	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}
	return value, kind
}

// ToBool 转为布尔类型
func ToBool(b any) (bool, error) {
	if xcore.IsZero(b) {
		return false, xerror.New("val is empty")
	}

	switch v := b.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return strconv.ParseBool(ToString(v))
	}
}

// ToSlice interface 转为 interface 切片
func ToSlice(val any) []any {
	if val == nil {
		return make([]interface{}, 0)
	}
	inst, kind := GetInst(val)
	if kind == reflect.Slice {
		v := make([]any, inst.Len())
		for i := 0; i < inst.Len(); i++ {
			v[i] = inst.Index(i).Interface()
		}
		return v

	} else {
		return []any{val}
	}
}

// WarpToMap 转换为 包装为 map转换函数
// 示例：s:[]string xarray.MapBy(s, xparse.WarpToMap(xparse.ToNumber[int64])) string[] 转为 int64[]
func WarpToMap[T, U any](parser func(item T) (U, error)) func(index int, item T) (U, error) {
	return func(index int, item T) (U, error) {
		return parser(item)
	}
}
