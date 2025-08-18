package x

import (
	"context"
	"os"
	"os/signal"

	"github.com/77d88/go-kit/plugins/xlog"
)

func Start() {
	if x.sf == nil {
		panic("server is nil please use Server")
	}
	x.wait.Wait()
	go func() {
		s, err := x.sf()
		if err != nil {
			panic(err)
		}
		x.Server = s
		if err != nil {
			panic(err)
		}
		s.Start()
	}()
	go func() {
		for _, f := range x.afterStart {
			go func(fn func()) {
				defer func() {
					if err := recover(); err != nil {
						xlog.Errorf(context.Background(), "after start panic: %v", err)
					}
				}()
				fn()
			}(f)
		}
	}()
	signal.Notify(x.QuitSignal, os.Interrupt)
	<-x.QuitSignal
	// 释放资源
	x.Close()
	return
}

func Close() {
	x.Close()
}

func Info() *XInfo {
	return x.Info
}

func CtxTraceId(ctx context.Context) (string, bool) {
	value := ctx.Value(X_TRACE_ID)
	s, ok := value.(string)
	return s, ok
}
