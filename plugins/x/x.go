package x

import (
	"os"
	"reflect"
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

func SetConfig(config *xconfig.Config) {
	x.Cfg = config
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

// FastInit 快速初始化构造器的参参数
func FastInit(construct interface{}) {
	typeOf := reflect.TypeOf(construct)
	if typeOf.Kind() == reflect.Func {
		for i := 0; i < typeOf.NumIn(); i++ {
			inType := typeOf.In(i)
			go func() {
				x.wait.Add(1)
				defer x.wait.Done()
				_, err := GetByType(inType)
				if err != nil {
					panic(err)
				}
			}()
		}
	}
}
