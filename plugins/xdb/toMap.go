package xdb

import (
	"context"
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"reflect"
	"strings"
	"time"
)

type states string

const ToMapIgnore states = "_"

type ToMapFunc func(value interface{}) (k string, v interface{}, err error)

// MapDateParse 将字符串转换为时间

func MapDateParse(field, layout string, defaultVal interface{}) ToMapFunc {
	return NewMapParse(field, func(value interface{}) (interface{}, error) {
		if value == nil || value.(string) == "" {
			return defaultVal, nil
		}
		return time.ParseInLocation(layout, value.(string), time.Local)
	})
}

// NewMapParse 自定义字段转换
func NewMapParse(field string, f func(value interface{}) (interface{}, error)) ToMapFunc {
	return func(value interface{}) (k string, v interface{}, err error) {
		v, err = f(value)
		return field, v, err
	}
}

// MapUpdateName 自定义字段名
func MapUpdateName(newName string) ToMapFunc {
	return func(value interface{}) (string, interface{}, error) {
		return newName, value, nil
	}
}

// MapUpdateId 更新ID
func MapUpdateId(id int64) func() int64 {
	return func() int64 {
		return id
	}
}

// toSqlMap 将对象转换为map mapping
//
//} 使用 NewMapParse 、MapUpdateName 、MapUpdateId 来创建 或者自己按下面规则创建
/**
mapping []map[string]interface{}{
	// a_field -> ToMapIgnore 忽略字段 ToMapIgnore
	"a_field": ToMapIgnore ,
	// 自定义转换 如果出现错误则不处理 方法返回的是 ToMapIgnore 则忽略字段
	"a_field": func(value interface{}) (interface{}, error) {}
	// 自定义转换 如果出现错误则不处理 方法返回的是 ToMapIgnore 则忽略字段 string 为自定义字段覆盖原始字段
	"a_field": func(value interface{}) (string,interface{}, error) {}
	// 自定义转换 类似于增加字段使用 方法返回的是 ToMapIgnore 则忽略字段
	"a_field": func() (interface{}, error) {}
	// 不管a_field 移除所有忽略大小写的ID 设置为返回值
	"a_field": func() int64
}
*/
func toSqlMap(c context.Context, obj interface{}, mapping ...interface{}) map[string]interface{} {
	value := reflect.ValueOf(obj)
	kind := value.Kind()
	if kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}

	result := make(map[string]interface{})
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)
		fieldName := fieldType.Name
		fValue := field.Interface()
		result[fieldName] = fValue
	}

	// 合并mapping  解析字段
	var parses = make(map[string]interface{})
	for i, m := range mapping {
		switch f := m.(type) {
		case map[string]interface{}:
			for k, v := range f {
				if f, ok := v.(func(value interface{}) (interface{}, error)); ok {
					parses[k] = ToMapFunc(func(value interface{}) (string, interface{}, error) {
						v, err := f(value)
						return k, v, err
					})
				} else if f, ok := v.(func() (interface{}, error)); ok {
					parses[k] = ToMapFunc(func(value interface{}) (string, interface{}, error) {
						v, err := f()
						return k, v, err
					})
				} else {
					parses[k] = v
				}
			}
		default:
			parses[fmt.Sprintf("key:%d", i)] = m
		}
	}

	for k, v := range parses {
		switch v := v.(type) {
		case func() int64:
			for k := range result {
				if strings.ToLower(k) == "id" {
					delete(result, k)
				}
			}
			result["id"] = v()
		case states:
			if v == ToMapIgnore { // 忽略字段
				delete(result, k)
				xlog.Tracef(c, "toSqlMap ignore field:%s", k)
			}
		case ToMapFunc:
			if newKey, newVal, err := v(result[k]); err == nil {
				if s, ok := newVal.(states); ok {
					if s == ToMapIgnore { // 方法返回的是_字符串则不处理
						delete(result, k)
						xlog.Tracef(c, "toSqlMap ignore field:%s", k)
						continue
					}

				} else {
					if k != newKey { // 移除旧key
						delete(result, k)
					}
					result[newKey] = newVal // 新key赋值
				}
			}
		default:
			if strings.ToLower(k) == "id" { // id 特殊处理
				for k := range result {
					if strings.ToLower(k) == "id" {
						delete(result, k)
					}
				}
			}
			result[k] = v
		}
	}

	return result
}
