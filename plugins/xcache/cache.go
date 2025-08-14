package xcache

import (
	"time"
)

type Cached interface {
	Set(key string, val any, expire time.Duration)                                                 // 缓存 nil一样缓存只要不报错
	Get(key string) (interface{}, bool)                                                            // 获取缓存 nil一样缓存只要不报错
	Del(key string)                                                                                // 删除缓存 nil一样缓存只要不报错
	Once(key string, expire time.Duration, fci00 func() (interface{}, error)) (interface{}, error) // 缓存一次 并扫描到val nil一样缓存只要不报错
}

// d 默认使用本地缓存
var d = New(5*time.Minute, 10*time.Minute)

// UseCache 自定义缓存
func UseCache(cache Cached) {
	d = cache
}

func Set(key string, value interface{}, expire time.Duration) {
	d.Set(key, value, expire)
}

func Get(key string) (interface{}, bool) {
	return d.Get(key)
}

func Del(key string) {
	d.Del(key)
}

func Once(key string, expire time.Duration, fc func() (interface{}, error)) (interface{}, error) {
	return d.Once(key, expire, fc)
}
