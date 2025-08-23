package xpg

import (
	"fmt"
	"reflect"

	"github.com/77d88/go-kit/basic/xstr"
	"github.com/go-viper/mapstructure/v2"
)

func Scan(list []map[string]any, dest any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	v := reflect.ValueOf(dest).Elem() // 获取指针指向的值
	isSlice := true
	var structType reflect.Type
	// 如果是切片指针，则查询多个记录
	if t.Elem().Kind() != reflect.Slice {
		isSlice = false
		structType = t.Elem()
	} else {
		structType = t.Elem().Elem()
	}

	// 处理查询结果
	if isSlice {
		var isPtr = false
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem() // *[]*struct
			isPtr = true
		}
		if structType.Kind() != reflect.Struct {
			return fmt.Errorf("result must be a pointer to a slice of structs")
		}
		// 创建切片并填充数据
		slice := reflect.MakeSlice(t.Elem(), 0, len(list))
		for _, data := range list {
			val, err := mapDecode(structType, data)
			if err != nil {
				return err
			}
			if isPtr {
				slice = reflect.Append(slice, val)
			} else {

				slice = reflect.Append(slice, val.Elem())
			}
		}
		v.Set(slice) // 将结果赋值给传入的参数
	} else {
		// 处理单个对象
		if len(list) > 0 {
			val, err := mapDecode(structType, list[0])
			if err != nil {
				return err
			}
			v.Set(val.Elem()) // 将结果赋值给传入的参数
		}
	}
	return nil
}

// mapDecode 将map数据映射到结构体
func mapDecode(t reflect.Type, data map[string]any) (reflect.Value, error) {
	// 现将map key转为camelCase

	for k, v := range data {
		camelCaseKey := xstr.CamelCase(k)
		data[camelCaseKey] = v
		//delete(data, k) // 不用删除旧的 key 满足使用db tag的映射
	}

	instance := reflect.New(t).Interface()
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &instance,
		WeaklyTypedInput: true, // 允许弱类型转换（如 map → struct）
		Metadata:         nil,
		TagName:          "db",
		Squash:           true,
	})
	if err != nil {
		return reflect.ValueOf(instance), err
	}
	err = decoder.Decode(data)
	if err != nil {
		return reflect.ValueOf(instance), err
	}
	return reflect.ValueOf(instance), nil
}

// check

// extractDBFields 递归提取结构体中的db字段
func extractDBFields(t reflect.Type) []string {
	var fields []string

	// 处理指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 只处理结构体类型
	if t.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过被忽略的字段
		if field.Tag.Get("db") == "-" {
			continue
		}

		// 处理匿名（嵌入）字段
		if field.Anonymous {
			// 递归获取嵌入结构体的字段
			embeddedFields := extractDBFields(field.Type)
			fields = append(fields, embeddedFields...)
			continue
		}

		// 获取db标签，如果没有则使用蛇形命名
		fieldTag := field.Tag.Get("db")
		if fieldTag == "" {
			fieldTag = xstr.SnakeCase(field.Name)
		}

		fields = append(fields, `"`+fieldTag+`"`)
	}

	return fields
}

// extractDBObj 递归提取结构体中的db字段 字段名默认 SnakeCase蛇形命名
func extractDBObj(obj interface{}) (map[string]interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj is nil")
	}

	value := reflect.ValueOf(obj)

	// 处理指针
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("obj is nil")
		}
		value = value.Elem()
	}

	// 只处理结构体
	if value.Kind() != reflect.Struct {
		// 本身就是map 则返回
		if value.Kind() == reflect.Map {
			return obj.(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("obj is not a struct")
	}

	result := make(map[string]interface{})
	extractDBFieldsFromValue(value, result)
	return result, nil
}

// extractDBFieldsFromValue 从 reflect.Value 中递归提取字段到 map
func extractDBFieldsFromValue(value reflect.Value, result map[string]interface{}) {
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// 跳过未导出字段
		if !fieldValue.CanInterface() {
			continue
		}

		// 跳过被忽略的字段
		if field.Tag.Get("db") == "-" {
			continue
		}

		// 处理匿名（嵌入）字段
		if field.Anonymous {
			// 递归处理嵌入结构体
			if fieldValue.Kind() == reflect.Struct {
				extractDBFieldsFromValue(fieldValue, result)
			} else if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
				if fieldValue.Elem().Kind() == reflect.Struct {
					extractDBFieldsFromValue(fieldValue.Elem(), result)
				}
			}
			continue
		}

		// 获取字段值
		if !fieldValue.IsValid() || !fieldValue.CanInterface() {
			continue
		}

		// 获取db标签，如果没有则使用蛇形命名
		fieldTag := field.Tag.Get("db")
		if fieldTag == "" {
			fieldTag = xstr.SnakeCase(field.Name)
		}

		// 直接赋值，如果有命名相同的字段会直接替换
		result[fieldTag] = fieldValue.Interface()
	}
}
