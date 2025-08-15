package xctx

import (
	"context"
	"sync"
)

const ContextKey = "_xctx/contextkey"

type Context interface {
	context.Context
	Set(key string, val any)
	Get(key string) (any, bool)
}

type Ctx struct {
	context.Context
	vals map[string]any
	// This mutex protects Keys map.
	mu sync.RWMutex
}

func (c *Ctx) Set(key string, val any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.vals == nil {
		c.vals = make(map[string]any)
	}

	c.vals[key] = val
}

func With(ctx context.Context) Context {
	return &Ctx{
		Context: ctx,
	}
}

func (c *Ctx) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.vals[key]
	return
}

func (c *Ctx) Value(key any) any {
	if key == ContextKey {
		return c
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	return c.Context.Value(key)
}

func (c *Ctx) Ctx() context.Context {
	return c
}

func WithVal(ctx context.Context, key string, val any) Context {
	c := &Ctx{
		Context: ctx,
	}
	c.Set(key, val)
	return c
}

func SetVal(c context.Context, key string, val any) context.Context {
	// 如果c 实现了 Context 接口 则直接调用
	if v, ok := c.(Context); ok {
		v.Set(key, val)
		return c
	}
	return context.WithValue(c, key, val) // 否则使用 WithValue
}
