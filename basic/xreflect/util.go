package xreflect

import (
	"fmt"
	"reflect"
)

func GetFieldVal[T any](obj T, field string) any {
	fields := GetAllFields(obj)
	f := fields[field]
	if !f.IsValid() {
		return nil
	}
	return f.GetVal()
}

// IsSlice 是否是切片
func IsSlice[T any](obj T) bool {
	_, kind := GetInst(obj)
	return kind == reflect.Slice
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

func IsPointer(val interface{}) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Ptr
}

// IsMap 是否是map
func IsMap[T any](obj T) bool {
	_, kind := GetInst(obj)
	return kind == reflect.Map
}

// ImplementsInterface 检查对象是否实现了指定接口
// obj: 要检查的对象
// ifacePtr: 接口类型的指针，例如 (*SomeInterface)(nil)
// 返回: 如果对象实现了接口则返回true，否则返回false
func ImplementsInterface(obj interface{}, ifacePtr interface{}) bool {
	if ifacePtr == nil {
		return false
	}

	// 获取接口的类型
	ifaceType := reflect.TypeOf(ifacePtr).Elem()

	// 检查ifacePtr是否是指向接口的指针
	if ifaceType.Kind() != reflect.Interface {
		return false
	}

	// 获取对象的类型
	objType := reflect.TypeOf(obj)

	// 如果obj是指针，需要检查指针类型是否实现了接口
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	// 检查类型是否实现了接口
	return reflect.PointerTo(objType).Implements(ifaceType)
}

// ErrorNoImplementsInterface  未实现接口的错误信息
var ErrorNoImplementsInterface = fmt.Errorf("object does not implement the interface")

// CallInterfaceMethod 动态调用对象的接口方法
// obj: 要调用方法的对象
// ifacePtr: 接口类型的指针，例如 (*SomeInterface)(nil)
// methodName: 要调用的方法名
// args: 方法参数(可选)
// 返回: 方法返回值(interface{}类型切片)和错误信息
func CallInterfaceMethod(obj interface{}, ifacePtr interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	// 1. 检查对象是否实现了接口
	if !ImplementsInterface(obj, ifacePtr) {
		return nil, ErrorNoImplementsInterface
	}

	// 2. 获取反射值
	val := reflect.ValueOf(obj)

	// 4. 查找方法
	method := val.MethodByName(methodName)
	if !method.IsValid() {

		// 如果是指针类型，再找一次值的方法
		if val.Kind() == reflect.Ptr {
			method = val.Elem().FieldByName(methodName)
		}
		if !method.IsValid() {
			return nil, fmt.Errorf("method %s not found", methodName)
		}
	}

	// 5. 准备参数
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 6. 调用方法
	results := method.Call(in)

	// 7. 转换返回值
	out := make([]interface{}, len(results))
	for i, result := range results {
		out[i] = result.Interface()
	}

	return out, nil
}

// IsNil 判断是否为interface nil
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
