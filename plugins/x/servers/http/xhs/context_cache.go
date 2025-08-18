package xhs

import (
	"context"
	"sync"

	"github.com/77d88/go-kit/plugins/xlog"

	"github.com/jinzhu/copier"
)

type CacheObj struct {
	value any
}

type ApiCache struct {
	cache map[string]interface{}

	mu sync.RWMutex
}

func (c *ApiCache) CacheOnce(key string, fci00 func() (interface{}, error)) (interface{}, error) {
	value, exists := c.CacheGet(key)
	if exists {
		return value, nil
	}
	f, err := fci00()
	if err != nil {
		return nil, err
	}
	c.Cache(key, f)
	return f, nil
}

func (c *ApiCache) CacheScanOnce(key string, val any, fci00 func() (interface{}, error)) error {
	once, err := c.CacheOnce(key, fci00)
	if err != nil {
		return err
	}
	if once == nil { // nil 不复制
		return nil
	}
	return copier.Copy(val, once)
}

func (c *ApiCache) Cache(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]any)
	}

	c.cache[key] = &CacheObj{value}
}

func (c *ApiCache) CacheGet(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.cache[key]
	if exists {
		return value.(*CacheObj).value, true
	}
	return nil, false
}
func (c *ApiCache) CacheScan(key string, val any) error {
	value, exists := c.CacheGet(key)
	if exists {
		if val != nil {
			err := copier.Copy(val, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *ApiCache) CacheDel(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		return
	}
	delete(c.cache, key)
}

func (c *ApiCache) CopyCacheMap() map[string]interface{} {
	clone := make(map[string]interface{}, len(c.cache))
	for key, value := range c.cache {
		clone[key] = value
	}
	return clone
}

func CacheOnce(c context.Context, key string, val any, fci00 func() (interface{}, error)) error {
	if c, ok := c.(*Ctx); ok {
		v, err := CacheOnceFunc(c, key, fci00)
		if err != nil {
			return err
		}
		if v == nil { // nil 不复制
			return nil
		}
		return copier.Copy(val, v)
	} else {
		xlog.Errorf(c, "context is not xapi.Ctx no support context cache")
		i, err := fci00()
		if err != nil {
			return err
		}
		return copier.Copy(val, i)
	}
}

func CacheOnceFunc(c context.Context, key string, fci00 func() (interface{}, error)) (interface{}, error) {
	if c, ok := c.(*Ctx); ok {
		value, exists := cacheGet(c, key)
		if exists {
			return value, nil
		}
		f, err := fci00()
		if err != nil {
			return nil, err
		}
		cacheSet(c, key, f)
		return f, nil
	} else {
		xlog.Errorf(c, "context is not xapi.Ctx no support context cache")
		return fci00()
	}
}

func cacheSet(c context.Context, key string, value any) {
	if c, ok := c.(*Ctx); ok {
		c.Set(key, &CacheObj{value})
	} else {
		xlog.Errorf(c, "context is not xapi.Ctx no support context cache")
	}
}

func cacheGet(c context.Context, key string) (value any, exists bool) {
	if c, ok := c.(*Ctx); ok {
		value, exists = c.Get(key)
		if exists {
			return value.(*CacheObj).value, true
		}
		return nil, false
	} else {
		xlog.Errorf(c, "context is not xapi.Ctx no support context cache")
		return nil, false
	}

}
