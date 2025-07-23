package xcache

import (
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"time"
)

type Cached interface {
	WarpKey(key string) string
	WarpValue(val interface{}) *Value
	UnWarpValue(val *Value) (interface{}, bool)
	Set(key string, val any, expire time.Duration)                                           // 缓存 nil一样缓存只要不报错
	Get(key string) (interface{}, bool)                                                      // 获取缓存 nil一样缓存只要不报错
	Del(key string)                                                                          // 删除缓存 nil一样缓存只要不报错
	Scan(key string, val any)                                                                // 缓存key扫描到val nil一样缓存只要不报错
	Once(key string, val any, expire time.Duration, fci00 func() (interface{}, error)) error // 缓存一次 并扫描到val nil一样缓存只要不报错
}

type Value struct {
	Val interface{}
}

// D 默认使用本地缓存
var D *Cached = xcore.V2p(NewLocalCache("globalCache", 5*time.Minute, 10*time.Minute))

// UseCache 自定义缓存
func UseCache(cache Cached) {
	D = &cache
}

type Namespace struct {
	Namespace string
}

func (l *Namespace) WarpKey(key string) string {
	return fmt.Sprintf("%s:%s", l.Namespace, key)
}

func (l *Namespace) WarpValue(val interface{}) *Value {
	return &Value{
		Val: val,
	}
}
func (l *Namespace) UnWarpValue(val *Value) (interface{}, bool) {
	if val == nil {
		return nil, false
	}
	return val.Val, true
}

func Set(key string, value interface{}, expire time.Duration) {
	if D == nil {
		return
	}
	cached := *D
	cached.Set(cached.WarpKey(key), cached.WarpValue(value), expire)
}

func Get(key string) (interface{}, bool) {
	if D == nil {
		return nil, false
	}
	cached := *D
	get, b := cached.Get(cached.WarpKey(key))
	if b {
		return cached.UnWarpValue(get.(*Value))
	}
	return get, b
}

func Del(key string) {
	if D == nil {
		return
	}
	cached := *D
	cached.Del(cached.WarpKey(key))
}

func Once(key string, val any, expire time.Duration, fci00 func() (interface{}, error)) error {
	if D == nil {
		return nil
	}
	cached := *D
	return cached.Once(key, val, expire, fci00)
}
