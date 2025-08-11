package x

import (
	"github.com/77d88/go-kit/plugins/xlog"
	"os"
	"os/signal"
	"reflect"
)

func Start() {
	Init()
	if x.sf == nil {
		panic("server is nil please use UseServer")
	}
	x.wait.Wait()
	go func() {
		s, err := x.sf()
		if err != nil {
			panic(err)
		}
		err = x.provide(func() EngineServer {
			return s
		})
		x.Server = s
		if err != nil {
			panic(err)
		}
		s.Start()
	}()
	signal.Notify(x.QuitSignal, os.Interrupt)
	<-x.QuitSignal

	// 关闭服务
	_ = x.invoke(func(s EngineServer) {
		s.Shutdown()
	})
	// 释放资源
	x.Close()
	return
}


func Use(constructor interface{}, delay ...bool) {
	Init()

	if reflect.TypeOf(constructor).Kind() == reflect.Func {
		// 判断这个构造函数的第一个返回值是否是 func() (EngineServer, error)
		switch f := constructor.(type) {
		case func() (EngineServer, error):
			x.sf = f
			return
		case func() EngineServer:
			x.sf = func() (EngineServer, error) {
				return f(), nil
			}
			return
		}
	} else {
		return
	}

	err := x.provide(constructor)
	if err != nil {
		panic(err)
	}
	if len(delay) > 0 && delay[0] {
		return
	}

	x.wait.Add(1)
	go func(a interface{}) {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf(nil, "use init panic: %v", err)
			}
			x.wait.Done()
		}()
		methodValue := reflect.ValueOf(a)
		methodType := methodValue.Type()
		allTypes := make([]reflect.Type, 0, methodType.NumOut())
		if methodType.Kind() == reflect.Func {
			for i := 0; i < methodType.NumOut(); i++ {
				allTypes = append(allTypes, methodType.Out(i))
			}
		}
		_, err := x.getInst(allTypes...)
		if err != nil {
			xlog.Fatalf(nil, "engine use error: %v", err)
		}
	}(constructor)
}

func Get[T any]() (*T, error) {
	var result T
	err := x.invoke(func(r T) {
		result = r
	})
	return &result, err
}
func Find(constructor interface{}) error {
	return x.invoke(constructor)
}

func Close() {
	x.Close()
}

func Info() *XInfo {
	return x.Info
}



func Config[T any](key string) (*T, error) {
	Init()
	var result T
	err := x.Cfg.ScanKey(key, &result)
	return &result, err
}

func ConfigString(key string) string {
	return x.Cfg.GetString(key)
}

func ConfigStringSlice(key string) []string {
	return x.Cfg.GetStringSlice(key)
}
func ConfigInt(key string) int {
	return x.Cfg.GetInt(key)
}
func ConfigIntSlice(key string) []int {
	return x.Cfg.GetIntSlice(key)
}
func ConfigBool(key string) bool {
	return x.Cfg.GetBool(key)
}
