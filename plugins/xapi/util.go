package xapi

import (
	"fmt"
	"github.com/77d88/go-kit/plugins/xe"
)

// GetByName 按名称获取实例(需配合命名注册使用)
func GetByName[T any](e *xe.Engine, name string) (T, error) {
	var result T
	var err error
	err2 := e.Container.Invoke(func(deps map[string]interface{}) {
		if val, ok := deps[name]; ok {
			if v, ok := val.(T); ok {
				result = v
			} else {
				err = fmt.Errorf("type mismatch for name: %s", name)
			}
		} else {
			err = fmt.Errorf("not found: %s", name)
		}
	})
	if err2 != nil {
		return result, err2
	}
	return result, err
}

// Get 通过类型获取容器中的实例
func Get[T any](e *xe.Engine) (T, error) {
	var result T
	err := e.Container.Invoke(func(val T) {
		result = val
	})
	return result, err
}
