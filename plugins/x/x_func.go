package x

import (
	"os"
	"os/signal"
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

func Config[T any](key string) (*T, error) {
	var result T
	err := x.Cfg.ScanKey(key, &result)
	return &result, err
}

func ConfigString(key string, defaultValue ...string) string {
	str := x.Cfg.GetString(key)
	if str == "" {
		if len(defaultValue) > 0 {
			str = defaultValue[0]
		}
	}
	return str
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
