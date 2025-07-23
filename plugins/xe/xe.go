package xe

import (
	"github.com/77d88/go-kit/basic/xconfig"
	"go.uber.org/dig"
	"time"
)

type LifeCycleDispose interface {
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
	Container *dig.Container
	Info      *XInfo
	Server    EngineServer
	Cfg       *xconfig.Config
}

var E *Engine // 持有一个服务实例

func New(cfg *xconfig.Config) *Engine {
	Container := dig.New()
	x := &Engine{
		Container: Container,
		Info: &XInfo{
			StartTime: time.Now(),
		},
		Cfg: cfg,
	}
	x.MustProvide(func() *Engine { return x })

	E = x
	return x
}

func (e *Engine) Provide(a interface{}, options ...dig.ProvideOption) error {
	return e.Container.Provide(a, options...)
}

func (e *Engine) Invoke(a interface{}, options ...dig.InvokeOption) error {
	return e.Container.Invoke(a, options...)
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
	e.MustInvoke(func(s EngineServer) {
		e.Server = s
		s.Start()
	})
	return
}

func (e *Engine) RunTime() time.Duration {
	return time.Since(e.Info.StartTime)
}

func (e *Engine) Dispose() error {
	var err error
	e.Container.Invoke(func(disposables ...LifeCycleDispose) {
		for _, disposable := range disposables {
			if disposable != nil {
				if err1 := disposable.Dispose(); err1 != nil {
					err = err1
				}
			}
		}
	})
	return err
}
