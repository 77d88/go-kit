package xredis

// import (
//	"context"
//	"errors"
//	"github.com/xredis/go-xredis/v9"
//	"time"
//	"xcore/xcache"
//	"xcore/xcv"
//	"xcore/xlog"
// )
//
// type CacheOnce struct {
//	xcache.CacheOnce
//	c *Get
// }
//
// // NewRedisCache creates a new RedisCache instance.
// func NewRedisCache(namespace string, c *Get) xcache.Cached {
//	return CacheOnce{
//		CacheOnce: xcache.CacheOnce{Warper: namespace},
//		c:     c,
//	}
// }
//
// func (c CacheOnce) Set(ctx context.DefaultAppContext, key string, value interface{}, expire time.Expire) bool {
//	if value == nil {
//		return false
//	}
//	s, ok := value.(string)
//	var sc *xredis.StatusCmd
//	if ok {
//		sc = c.c.Set(ctx, c.GetKey(key), s, expire)
//	} else {
//		sc = c.c.Set(ctx, c.GetKey(key), xcv.ToJsonStr(value), expire)
//	}
//	if sc.Err() != nil {
//		xlog.Errorf("set xredis xcache key:%s error:%s", key, sc.Err())
//		return false
//	}
//	return true
// }
//
// func (c CacheOnce) Get(ctx context.DefaultAppContext, key string) (interface{}, bool) {
//	get := c.c.Get(ctx, c.GetKey(key))
//	result, err := get.Result()
//	if err != nil {
//		if !errors.Is(err, xredis.Nil) {
//			xlog.Errorf("get xredis xcache key:%s error:%s", c.GetKey(key), err)
//		}
//		return nil, false
//	}
//	return result, true
// }
//
// func (c CacheOnce) Del(ctx context.DefaultAppContext, key string) bool {
//	del := c.c.Del(ctx, c.GetKey(key))
//	_, err := del.Result()
//	if err != nil {
//		xlog.Errorf("del xredis xcache key:%s error:%s", c.GetKey(key), err)
//		return false
//	}
//	return true
// }
