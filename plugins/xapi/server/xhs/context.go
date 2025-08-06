package xhs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
)

const ctxThisKey = "_xe/CTX_KEY_CONTEXT"

type Ctx struct {
	*gin.Context
	ContextAuth
	ApiCache
	Result     interface{}
	PrintStack bool
	TraceId    int64
	test       bool
	Server     *HttpServer
}

type CopyContext struct {
	context.Context
	ContextAuth
	ApiCache
	TraceId int64
}

func newCtx(c *gin.Context, x *HttpServer) *Ctx {
	if value, exists := c.Get(ctxThisKey); exists {
		if c2, ok := value.(*Ctx); ok {
			return c2
		}
	}
	d := &Ctx{
		Context:    c,
		PrintStack: true,
		Result:     nil,
		TraceId:    xid.NextId(),
		Server:     x,
	}
	c.Set(ctxThisKey, d)
	return d
}

func NewTestContext() *Ctx {
	return &Ctx{
		test:    true,
		Context: &gin.Context{},
	}
}

func (c *Ctx) Value(key any) any {
	if key == ctxThisKey {
		return c
	}
	if key == xlog.CtxLogParam {
		return logFields(c)
	}
	return c.Context.Value(key)
}

func (c *Ctx) Set(key string, val any) {
	c.Context.Set(key, val)
}

func (c *Ctx) Send(v interface{}) {
	c.Result = NewResp(v)

	if c.test {
		indent, err := json.MarshalIndent(c.Result, "", "  ")
		if err != nil {
			fmt.Printf("test result:\n%v\n", c.Result)
		} else {
			fmt.Printf("test result:\n%v\n", string(indent))
		}
	}
}

func (c *Ctx) SendError(err interface{}) {
	e := xerror.New(err)
	c.Result = e
	_ = c.Error(e)
}

// Copy 拷贝最终上下文 用于传输获其他现场调用
func (c *Ctx) Copy() context.Context {
	return CopyContext{
		Context:     context.WithValue(c.Context.Copy(), xlog.CtxLogParam, logFields(c)),
		ContextAuth: c.ContextAuth,
		ApiCache: ApiCache{
			cache: c.CopyCacheMap(),
		},
		TraceId: c.TraceId,
	}
}
