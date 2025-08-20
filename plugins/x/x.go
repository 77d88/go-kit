package x

import (
	"os"
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/xlog"
)

const (
	X_TRACE_ID = "__%X-Trace-Id"
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
	Info       *XInfo
	Cfg        *xconfig.Config
	registry   []interface{}
	QuitSignal chan os.Signal
	wait       sync.WaitGroup
	Server     EngineServer
	sf         func() (EngineServer, error)
	dises      []Disposer
	afterStart []func()
}

var x *Engine // 持有一个服务实例

func init() {
	engine := &Engine{
		Info: &XInfo{
			StartTime: time.Now(),
		},
		Cfg:        nil, // 这个通过 use 注入 或者setConfig
		QuitSignal: make(chan os.Signal),
	}
	container.UseInitAfter(func(key string, value interface{}, fc bool) {

		if v, ok := value.(Disposer); ok {
			x.dises = append(x.dises, v)
		}

		if v, ok := value.(*xconfig.Config); ok {
			x.Cfg = v
		}

		if v, ok := value.(EngineServer); ok {
			x.sf = func() (EngineServer, error) {
				return v, nil
			}
		}

		if v, ok := value.(func() (EngineServer, error)); ok {
			x.sf = v
		}

		if v, ok := value.(func() EngineServer); ok {
			x.sf = func() (EngineServer, error) {
				return v(), nil
			}
		}

	})
	Use(engine)
	x = engine
}

func (e *Engine) Must(constructorOrValue interface{}, name ...string) {
	key := Use(constructorOrValue, name...)
	e.wait.Go(func() {
		_, err := Get[any](key)
		if err != nil {
			panic(err)
		}
	})
}

func (e *Engine) AfterStart(constructor func()) {
	e.afterStart = append(e.afterStart, constructor)
}

func Must(constructorOrValue interface{}, name ...string) {
	x.Must(constructorOrValue, name...)
}

// AfterStart 服务启动后执行 所有must执行完毕
func AfterStart(constructor func()) {
	x.AfterStart(constructor)
}

func (e *Engine) Close() {
	e.Server.Shutdown()
	xlog.Warnf(nil, "Engine close")
	for _, inst := range e.dises {
		if err1 := inst.Dispose(); err1 != nil {
			xlog.Errorf(nil, "dispose error: %v", err1)
		}
	}
}
