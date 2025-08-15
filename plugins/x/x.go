package x

import (
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/xlog"
	"go.uber.org/dig"
)

type Disposer interface {
	Dispose() error // 释放资源
}

type EngineServer interface {
	Start()
	Shutdown()
	Name() string
}
type XInfo struct {
	StartTime time.Time
}

type Engine struct {
	Container  *dig.Container
	Info       *XInfo
	Cfg        *xconfig.Config
	registry   []interface{}
	QuitSignal chan os.Signal
	wait       sync.WaitGroup
	Server     EngineServer
	sf         func() (EngineServer, error)
}

var x *Engine // 持有一个服务实例

func init() {
	engine := &Engine{
		Container: dig.New(),
		Info: &XInfo{
			StartTime: time.Now(),
		},
		Cfg:        nil, // 这个通过 use 注入 或者setConfig
		QuitSignal: make(chan os.Signal),
	}
	err := engine.provide(func() *Engine { return engine })
	if err != nil {
		panic(err)
	}
	x = engine
}

func SetConfig(config *xconfig.Config) {
	x.Cfg = config
	err := x.provide(func() *xconfig.Config { return x.Cfg })
	if err != nil {
		panic(err)
	}
}

func (e *Engine) provide(constructor interface{}, options ...dig.ProvideOption) error {
	err := e.Container.Provide(constructor, options...)
	if err == nil {
		e.registry = append(e.registry, constructor) // 记录一下构造函数
	}
	return err
}

func (e *Engine) invoke(constructor interface{}, options ...dig.InvokeOption) error {
	err := e.Container.Invoke(constructor, options...)
	return err
}

func (e *Engine) RunTime() time.Duration {
	return time.Since(e.Info.StartTime)
}

func (e *Engine) getAllUseType() []reflect.Type {
	allTypes := make([]reflect.Type, 0, len(e.registry))
	for _, v := range e.registry {
		methodValue := reflect.ValueOf(v)
		methodType := methodValue.Type()
		if methodType.Kind() == reflect.Func {
			for i := 0; i < methodType.NumOut(); i++ {
				allTypes = append(allTypes, methodType.Out(i))
			}
		}

	}
	// 去重
	return xarray.UniqueBy(allTypes, func(item reflect.Type) string {
		return item.String()
	})
}

func (e *Engine) Close() {
	xlog.Warnf(nil, "Engine close")
	getInst, _ := e.getInst(e.getAllUseType()...)
	for _, inst := range getInst {
		if disposable, ok := inst.(Disposer); ok {
			if disposable != nil {
				if err1 := disposable.Dispose(); err1 != nil {
					xlog.Errorf(nil, "dispose error: %v", err1)
				}
			}
		}
	}
}

func (e *Engine) getInst(types ...reflect.Type) ([]interface{}, error) {
	value, err := e.getInstValue(types...)
	if err != nil {
		return nil, err
	}
	return xarray.Map(value, func(index int, item reflect.Value) interface{} {
		return item.Interface()
	}), nil
}
func GetInst(types ...reflect.Type) ([]interface{}, error) {
	return x.getInst(types...)
}

func (e *Engine) getInstValue(types ...reflect.Type) ([]reflect.Value, error) {
	allInst := make([]reflect.Value, 0, len(types))
	injectFuncType := reflect.FuncOf(types, nil, false)
	invokeFunc := reflect.MakeFunc(injectFuncType, func(args []reflect.Value) []reflect.Value {
		allInst = append(allInst, args...)
		return nil
	})
	if err := e.invoke(invokeFunc.Interface()); err != nil {
		xlog.Errorf(nil, "dig inject error: %v", err)
		return nil, err
	}
	return allInst, nil
}

func GetInstValue(types ...reflect.Type) ([]reflect.Value, error) {
	return x.getInstValue(types...)
}
