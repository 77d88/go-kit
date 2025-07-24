package xreflect

import (
	"reflect"
	//"github.com/patrickmn/go-cache"
)

type Field struct {
	reflect.StructField
	reflect.Value
}

func (f Field) GetVal() interface{} {
	return f.Interface()
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

// GetAllFields 获取所有字段
func GetAllFields[T any](obj T) map[string]Field {
	value, _ := GetInst(obj)
	vt := value.Type()
	numField := value.NumField()
	fields := make(map[string]Field, numField)
	for i := 0; i < numField; i++ {
		field := Field{
			StructField: vt.Field(i),
			Value:       value.Field(i),
		}
		fields[field.Name] = field
	}
	return fields
}

func SetFieldValue(obj any, field string, value any) {
	if obj == nil {
		return
	}
	if IsSlice(obj) {
		for _, v := range obj.([]any) {
			SetFieldValue(v, field, value)
		}
		return
	}
	inst, _ := GetInst(obj)
	fieldByName := inst.FieldByName(field)
	if !fieldByName.IsValid() {
		return
	}
	fieldByName.Set(reflect.ValueOf(value))

}
