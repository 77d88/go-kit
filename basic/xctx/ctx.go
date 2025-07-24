package xctx

import "context"

type CtxSet interface {
	Set(key string, val any)
}

func SetVal(c context.Context, key string, val any) context.Context {
	// 如果c 实现了 CtxSet 接口 则直接调用
	if v, ok := c.(CtxSet); ok {
		v.Set(key, val)
		return c
	}
	return context.WithValue(c, key, val) // 否则使用 WithValue
}
