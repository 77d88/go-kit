package xe

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/xlog"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"reflect"
	"time"
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
}

var E *Engine // 持有一个服务实例

func New(cfg *xconfig.Config) *Engine {
	Container := dig.New()
	x := &Engine{
		Container: Container,
		Info: &XInfo{
			StartTime: time.Now(),
		},
		Cfg:        cfg,
		QuitSignal: make(chan os.Signal),
	}
	x.MustProvide(func() *Engine { return x })

	E = x
	return x
}

func (e *Engine) Provide(a interface{}, options ...dig.ProvideOption) error {
	err := e.Container.Provide(a, options...)
	if err == nil {
		e.registry = append(e.registry, a) // 记录一下构造函数
	}
	return err
}

func (e *Engine) Invoke(a interface{}, options ...dig.InvokeOption) error {
	err := e.Container.Invoke(a, options...)
	return err
}

func Invoke(a interface{}, options ...dig.InvokeOption) error {
	return E.Invoke(a, options...)
}

func Provide(a interface{}, options ...dig.ProvideOption) error {
	return E.Provide(a, options...)
}

func MustProvide(a interface{}, options ...dig.ProvideOption) {
	E.MustProvide(a, options...)
}

func MustInvoke(a interface{}, options ...dig.InvokeOption) {
	E.MustInvoke(a, options...)
}

func (e *Engine) MustProvide(a interface{}, options ...dig.ProvideOption) *Engine {
	var err error
	err = e.Provide(a, options...)
	if err != nil {
		panic(err)
	}
	return e
}


func (e *Engine) MustInvoke(a interface{}, options ...dig.InvokeOption) *Engine {
	var err error
	err = e.Invoke(a, options...)
	if err != nil {
		panic(err)
	}
	return e
}

func (e *Engine) Start() {

	go func() {
		err := e.Invoke(func(s EngineServer) {
			s.Start()
		})
		if err != nil {
			xlog.Errorf(nil, "server start error error: %v", err)
		}
	}()
	signal.Notify(e.QuitSignal, os.Interrupt)
	<-e.QuitSignal

	// 关闭服务
	_ = e.Invoke(func(s EngineServer) {
		s.Shutdown()
	})
	// 释放资源
	e.Close()
	return
}

func (e *Engine) RunTime() time.Duration {
	return time.Since(e.Info.StartTime)
}

func (e *Engine) Close() {
	xlog.Warnf(nil, "Engine close")
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
	allTypes = xarray.UniqueBy(allTypes, func(item reflect.Type) string {
		return item.String()
	})
	getInst, _ := e.GetInst(allTypes...)
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

func (e *Engine) GetInst(types ...reflect.Type) ([]interface{}, error) {
	value, err := e.GetInstValue(types...)
	if err != nil {
		return nil, err
	}
	return xarray.Map(value, func(index int, item reflect.Value) interface{} {
		return item.Interface()
	}), nil
}

func (e *Engine) GetInstValue(types ...reflect.Type) ([]reflect.Value, error) {
	allInst := make([]reflect.Value, 0, len(types))
	injectFuncType := reflect.FuncOf(types, nil, false)
	invokeFunc := reflect.MakeFunc(injectFuncType, func(args []reflect.Value) []reflect.Value {
		allInst = append(allInst, args...)
		return nil
	})
	if err := e.Invoke(invokeFunc.Interface()); err != nil {
		xlog.Errorf(nil, "dig inject error: %v", err)
		return nil, err
	}
	return allInst, nil
}
