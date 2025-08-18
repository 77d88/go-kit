package x

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/77d88/go-kit/basic/xerror"
)

var typeKeyCache sync.Map
var container = NewContainer()

// Container 容器结构体
type Container struct {
	vals       map[string]interface{}
	funcs      map[string]interface{}
	building   map[string]*buildResult // 添加这行，用于跟踪正在构建的实例
	buildMutex sync.Mutex              // 添加这行，用于保护building map
	mu         sync.RWMutex
	bf         func(key string, value interface{}, fc bool)
}

func (c *Container) setVal(key string, value interface{}) {
	container.vals[key] = value
	if c.bf != nil {
		c.bf(key, value, false)
	}
}

func (c *Container) setFunc(key string, value interface{}) {
	container.funcs[key] = value
	if c.bf != nil {
		c.bf(key, value, true)
	}
}

func (c *Container) UseInitAfter(f func(key string, value interface{}, fc bool)) {
	c.bf = f
}

type buildResult struct {
	ready chan struct{} // 用于通知实例构建完成
	value interface{}   // 构建完成的值
	err   error         // 构建过程中的错误
}

// NewContainer 创建新容器
func NewContainer() *Container {
	return &Container{
		vals:     make(map[string]interface{}),
		funcs:    make(map[string]interface{}), // 添加这行
		building: make(map[string]*buildResult),
	}
}

// Use 往容器里面添加实例 可以是构造函数 也可以是直接的实例
func Use(constructorOrValue interface{}, name ...string) string {
	container.mu.Lock()
	defer container.mu.Unlock()

	typeOf := reflect.TypeOf(constructorOrValue)
	if typeOf.Kind() == reflect.Func {
		// 解析函数 函数只能返回最多两个参数 一个是实例 一个是错误
		numOut := typeOf.NumOut()
		if numOut == 0 || numOut > 2 {
			panic("constructor function must return either one value or two values (instance and error)")
		}
		// 获取第一个返回参数的类型作为键名
		returnType := typeOf.Out(0)
		// 如果有name那么使用name 如果不是 那么使用参数类型放入到components中
		if len(name) > 0 {
			key := name[0]
			if _, exists := container.funcs[key]; exists {
				panic(xerror.Newf("function with name %s already exists", key))
			}
			if _, exists := container.vals[key]; exists {
				panic(xerror.Newf("value with name %s already exists", key))
			}
			container.setFunc(key, constructorOrValue)
			return key
		} else {
			key := getTypeKey(returnType)
			if _, exists := container.vals[key]; exists {
				panic(xerror.Newf("value with type %s already exists", key))
			}
			if _, exists := container.funcs[key]; exists {
				panic(xerror.Newf("function with type %s already exists", key))

			}
			container.setFunc(key, constructorOrValue)
			return key
		}
	} else {
		// 处理直接的实例
		if len(name) > 0 {
			key := name[0]
			if _, exists := container.vals[key]; exists {
				panic(xerror.Newf("value with name %s already exists", key))
			}
			container.setVal(key, constructorOrValue)
			return key
		} else {
			// 通过获取实例的类型组成组件名
			key := getTypeKey(typeOf)
			if _, exists := container.vals[key]; exists {
				panic(xerror.Newf("value with type %s already exists", key))
			}
			container.setVal(key, constructorOrValue)
			return key
		}
	}
}

// Get 获取组件
func Get[T any](name ...string) (T, error) {
	var t T
	var key string
	if len(name) > 0 {
		key = name[0]
	} else {
		// 通过获取实例的类型组成组件名
		key = getTypeKey(reflect.TypeOf(t))
	}
	// 先尝试从vals中查找已经实例化的对象（使用读锁）
	container.mu.RLock()
	if val, ok := container.vals[key]; ok {
		container.mu.RUnlock()
		if result, ok := val.(T); ok {
			return result, nil
		}
		return t, xerror.Newf("%s 不支持类型 %s key Type is %s", val, key, getTypeKey(reflect.TypeOf(val)))
	}
	constructor, ok := container.funcs[key]
	container.mu.RUnlock()
	if ok {
		return resolveConstructor[T](constructor, key)
	}

	return t, fmt.Errorf("component %s not found", key)
}

func GetByType(t reflect.Type) (interface{}, error) {
	// 获取这个type的实例
	return Get[any](getTypeKey(t))
}

// Find 查找组件
// 通过反射解析结构体字段，并从容器中查找对应的依赖项进行注入
func Find[T any]() (T, error) {
	var t T
	tType := reflect.TypeOf(t)
	tValue := reflect.ValueOf(&t).Elem()

	// 检查是否为结构体或结构体指针
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}

	// 确保是结构体类型
	if tType.Kind() != reflect.Struct {
		return t, xerror.New("Find only supports struct types")
	}

	// 遍历结构体字段
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		fieldValue := tValue.Field(i)

		// 跳过不可设置的字段
		if !fieldValue.CanSet() {
			continue
		}

		// 解析tag中的name
		tag := field.Tag.Get("x")
		name := parseTagName(tag)

		// 获取字段类型对应的key
		var key string
		if name != "" {
			key = name
		} else {
			key = getTypeKey(field.Type)
		}

		// 从容器中获取对应的值
		get, err := Get[any](key)
		if err != nil {
			return t, xerror.Newf("failed to get value for field %s: %v", field.Name, err)
		}
		// 设置字段值
		if get != nil {
			valValue := reflect.ValueOf(get)
			if valValue.Type().AssignableTo(field.Type) {
				fieldValue.Set(valValue)
			} else {
				return t, xerror.Newf("cannot assign %s to field %s of type %s",
					valValue.Type(), field.Name, field.Type)
			}
		}
	}
	return t, nil
}

// parseTagName 解析tag中的name值
func parseTagName(tag string) string {
	if tag == "" {
		return ""
	}

	// 简单解析 x:"name=123" 格式
	// 实际项目中可能需要更复杂的解析逻辑
	for _, part := range strings.Split(tag, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "name=") {
			return strings.TrimPrefix(part, "name=")
		}
	}
	return ""
}

// resolveConstructor 解析构造函数并创建实例
func resolveConstructor[T any](constructor interface{}, key string) (T, error) {
	// 先检查是否已经有构建好的实例
	container.mu.RLock()
	if val, ok := container.vals[key]; ok {
		container.mu.RUnlock()
		if result, ok := val.(T); ok {
			return result, nil
		}
		var t T
		return t, xerror.Newf("type mismatch for key %s: expected %T, got %T", key, t, val)
	}
	container.mu.RUnlock()

	// 双重检查锁模式确保只有一个goroutine会构建实例
	container.buildMutex.Lock()
	// 再次检查是否已经有构建好的实例或者正在构建中
	container.mu.RLock()
	if val, ok := container.vals[key]; ok {
		container.mu.RUnlock()
		container.buildMutex.Unlock()
		if result, ok := val.(T); ok {
			return result, nil
		}
		var t T
		return t, xerror.Newf("type mismatch for key %s: expected %T, got %T", key, t, val)
	}

	// 检查是否正在构建中
	if building, ok := container.building[key]; ok {
		// 如果正在构建中，等待构建完成
		container.mu.RUnlock()
		container.buildMutex.Unlock()
		<-building.ready // 等待构建完成
		if building.err != nil {
			var t T
			return t, building.err
		}
		if result, ok := building.value.(T); ok {
			return result, nil
		}
		var t T
		return t, xerror.Newf("type mismatch for key %s: expected %T, got %T", key, t, building.value)
	}

	// 标记为正在构建中
	result := &buildResult{
		ready: make(chan struct{}),
	}
	container.building[key] = result
	container.mu.RUnlock()
	container.buildMutex.Unlock()

	// 执行实际的构建过程
	value, err := doResolveConstructor[T](constructor, key)

	// 构建完成后，通知所有等待者
	container.buildMutex.Lock()
	delete(container.building, key)
	result.value = value
	result.err = err
	close(result.ready)
	container.buildMutex.Unlock()

	return value, err
}

// doResolveConstructor 实际执行构造函数解析的逻辑
func doResolveConstructor[T any](constructor interface{}, key string) (T, error) {
	constructorValue := reflect.ValueOf(constructor)
	constructorType := constructorValue.Type()

	// 获取函数输入参数
	numIn := constructorType.NumIn()
	in := make([]reflect.Value, numIn)

	// 递归解析依赖项
	for i := 0; i < numIn; i++ {
		argType := constructorType.In(i)
		argKey := getTypeKey(argType)

		// 先从vals中查找已经实例化的对象
		container.mu.RLock()
		if val, ok := container.vals[argKey]; ok {
			container.mu.RUnlock()
			in[i] = reflect.ValueOf(val)
			continue
		}
		container.mu.RUnlock()

		// 再从funcs中查找构造函数
		container.mu.RLock()
		if funcVal, ok := container.funcs[argKey]; ok {
			container.mu.RUnlock()
			// 递归解析构造函数
			resolvedVal, err := resolveConstructor[T](funcVal, argKey)
			if err != nil {
				return *new(T), err
			}
			in[i] = reflect.ValueOf(resolvedVal)
			continue
		}
		container.mu.RUnlock()
		return *new(T), xerror.Newf("dependency %s not found for constructor %s", argKey, key)
	}

	// 调用构造函数
	results := constructorValue.Call(in)

	// 处理返回值
	i := results[0].Interface()
	if len(results) == 1 {
		// 确保在更新容器时使用写锁
		container.mu.Lock()
		container.setVal(key, i)
		container.mu.Unlock()
		return i.(T), nil
	} else if len(results) == 2 {
		instance := i
		errValue := results[1]
		if errValue.IsNil() {
			// 确保在更新容器时使用写锁
			container.mu.Lock()
			container.setVal(key, i)
			container.mu.Unlock()
			return instance.(T), nil
		}
		// 如果有错误返回值且不为nil，则返回该错误
		if err, ok := errValue.Interface().(error); ok {
			return instance.(T), err
		}
		return instance.(T), xerror.New("constructor returned non-nil error")
	}

	return *new(T), xerror.New("unexpected number of return values")
}

func getTypeKey(t reflect.Type) string {
	if t == nil {
		return ""
	}

	if key, ok := typeKeyCache.Load(t); ok {
		return key.(string)
	}

	// 处理指针类型
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	key := ""
	// 处理切片、数组、map等复合类型
	switch t.Kind() {
	case reflect.Slice:
		key = "[]" + getTypeKey(t.Elem())
	case reflect.Array:
		key = fmt.Sprintf("[%d]%s", t.Len(), getTypeKey(t.Elem()))
	case reflect.Map:
		key = fmt.Sprintf("map[%s]%s", getTypeKey(t.Key()), getTypeKey(t.Elem()))
	case reflect.Chan:
		key = fmt.Sprintf("chan %s", getTypeKey(t.Elem()))
	case reflect.Interface:
		if t.PkgPath() == "" {
			key = t.Name() // 内置接口类型如error
		}
		key = fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	case reflect.Struct:
		// 匿名结构体特殊处理
		if t.Name() == "" {
			key = fmt.Sprintf("%s.anonymous", t.PkgPath())
		}
		fallthrough
	default:
		// 基本类型和命名类型
		if t.PkgPath() == "" {
			key = t.Name() // 内置类型如int, string等
		}
		key = fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	}
	typeKeyCache.Store(t, key)
	return key
}
