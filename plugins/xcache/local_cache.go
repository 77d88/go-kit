package xcache

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xlog"
	"time"

	"github.com/jinzhu/copier"
	c "github.com/patrickmn/go-cache"
)

type LocalCache struct {
	*c.Cache
	Namespace
}

// NewLocalCache creates a new local xcache with the given namespace, default expiration, and cleanup interval.
func NewLocalCache(namespace string, defaultExpire, cleanupInterval time.Duration) Cached {
	return &LocalCache{
		Cache:     c.New(defaultExpire, cleanupInterval),
		Namespace: Namespace{Namespace: namespace},
	}
}


func (l *LocalCache) Set(key string, val any, expire time.Duration) {
	l.Cache.Set(l.WarpKey(key), val, expire)
}

func (l *LocalCache) Get(key string) (interface{}, bool) {
	return l.Cache.Get(l.WarpKey(key))
}

func (l *LocalCache) Del(key string) {
	l.Cache.Delete(l.WarpKey(key))
}

func (l *LocalCache) Scan(key string, val any) {
	get, b := l.Get(key)
	if !b {
		return
	}
	err := copier.Copy(val, get)
	xlog.Errorf(nil, "Scan %v %v %v", key, b, err)
}

func (l *LocalCache) Once(key string, val any, expire time.Duration, fci00 func() (interface{}, error)) error {
	get, b := l.Get(key)
	if !b {
		v, err := fci00()
		if err != nil {
			xlog.Errorf(nil, "xcache once %v %v", key, err)
			return xerror.Newf("xcache once %v %v", key, err)
		}
		l.Set(key, v, expire)
		get = v
	}
	if val != nil {
		return copier.Copy(val, get)
	}
	return nil
}
