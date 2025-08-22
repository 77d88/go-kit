package xcache

import (
	"fmt"
	"time"

	"github.com/77d88/go-kit/plugins/xlog"

	"github.com/patrickmn/go-cache"
)

type LocalCache struct {
	*cache.Cache
}

// New creates a new local xcache with the given namespace, default expiration, and cleanup interval.
func New(defaultExpire, cleanupInterval time.Duration) Cached {
	return &LocalCache{
		Cache: cache.New(defaultExpire, cleanupInterval),
	}
}

func (l *LocalCache) Set(key string, val any, expire time.Duration) {
	l.Cache.Set(key, val, expire)
}

func (l *LocalCache) Get(key string) (interface{}, bool) {
	return l.Cache.Get(key)
}

func (l *LocalCache) Del(key string) {
	l.Cache.Delete(key)
}

func (l *LocalCache) Once(key string, expire time.Duration, fci00 func() (interface{}, error)) (interface{}, error) {
	get, b := l.Get(key)
	if !b {
		v, err := fci00()
		if err != nil {
			xlog.Errorf(nil, "xcache once %v %v", key, err)
			return nil, fmt.Errorf("xcache once %v %w", key, err)
		}
		l.Set(key, v, expire)
		get = v
	}
	return get, nil
}
