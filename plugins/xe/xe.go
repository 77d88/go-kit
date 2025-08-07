package xe

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/xlog"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"reflect"
	"sync"
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
	wait       sync.WaitGroup
	Server     EngineServer
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
	err := x.Provide(func() *Engine { return x })
	if err != nil {
		panic(err)
	}

	E = x
	return x
}

// Use 添加依赖 这里是强制依赖立即初始化
func (e *Engine) Use(a interface{}) *Engine {
	err := Provide(a)
	if err != nil {
		panic(err)
		return e
	}
	e.wait.Add(1)
	go func(a interface{}) {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf(nil, "use init panic: %v", err)
			}
			e.wait.Done()
		}()
		methodValue := reflect.ValueOf(a)
		methodType := methodValue.Type()
		allTypes := make([]reflect.Type, 0, methodType.NumOut())
		if methodType.Kind() == reflect.Func {
			for i := 0; i < methodType.NumOut(); i++ {
				allTypes = append(allTypes, methodType.Out(i))
			}
		}
		_, err := e.GetInst(allTypes...)
		if err != nil {
			xlog.Fatalf(nil, "engine use error: %v", err)
		}
	}(a)
	return e
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

func (e *Engine) UseServer(f func(e *Engine) (EngineServer, error)) *Engine {
	s, err := f(e)
	if err != nil {
		panic(err)
	}
	err = e.Invoke(func() EngineServer {
		return s
	})
	e.Server = s
	if err != nil {
		panic(err)
	}
	return e
}

func (e *Engine) Start() {
	if e.Server == nil {
		panic("server is nil please use UseServer")
	}
	e.wait.Wait()
	go func() {
		e.Server.Start()
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
	getInst, _ := e.GetInst(e.getAllUseType()...)
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
