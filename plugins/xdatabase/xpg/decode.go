package xpg

import (
	"reflect"

	"github.com/77d88/go-kit/basic/xstr"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// mapDecode 将map数据映射到结构体
func mapDecode(t reflect.Type, data map[string]any) (reflect.Value, error) {
	// 现将map key转为camelCase
	instance := reflect.New(t).Interface()

	for k, v := range data {
		camelCaseKey := xstr.CamelCase(k)
		data[camelCaseKey] = v
		switch t := v.(type) {
		case pgtype.Numeric:
			value, err := t.Float64Value()
			if err != nil {
				return reflect.ValueOf(instance), err
			}
			data[camelCaseKey] = value.Float64
			data[k] = value.Float64
		}
		//delete(data, k) // 不用删除旧的 key 满足使用db tag的映射
	}

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
