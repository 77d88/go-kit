package xhs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
)

const ctxThisKey = "_xe/CTX_KEY_CONTEXT"

type Ctx struct {
	*gin.Context
	Auth *ContextAuth
	ApiCache
	Result     interface{}
	PrintStack bool
	TraceId    int64
	test       bool
	Server     *HttpServer
	errors     []*xerror.Error
}

type CopyContext struct {
	context.Context
	Auth *ContextAuth
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

func NewTestContext(data ...interface{}) *Ctx {
	c := &Ctx{
		test:    true,
		Context: &gin.Context{},
	}
	if len(data) > 0 {
		c.Set("testdata", data[0])
	}
	return c
}

func (c *Ctx) ShouldBind(obj any) error {
	if obj == nil {
		return nil
	}
	// 如果请求的body是空的，则返回错误
	if c.Request.ContentLength == 0 {
		return nil
	}
	if c.test {
		obj = c.Value("testdata")
	}
	return c.Context.ShouldBind(obj)
}

func ShouldBind[T any](c *Ctx) (T, error) {
	var t T
	err := c.ShouldBind(&t)
	return t, err
}

func (c *Ctx) Value(key any) any {
	if key == ctxThisKey {
		return c
	}
	if key == xlog.CtxLogParam {
		return logFields(c)
	}
	if key == x.X_TRACE_ID {
		return c.TraceId
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
	if c.errors == nil {
		c.errors = make([]*xerror.Error, 0)
	}
	c.errors = append(c.errors, e)
}

func (c *Ctx) GetError() error {
	if c.errors == nil {
		return nil
	}
	// 只获取最后一个
	return c.errors[len(c.errors)-1]
}

// Copy 拷贝最终上下文 用于传输获其他现场调用
func (c *Ctx) Copy() context.Context {
	return CopyContext{
		Context: context.WithValue(c.Context.Copy(), xlog.CtxLogParam, logFields(c)),
		Auth:    c.Auth,
		ApiCache: ApiCache{
			cache: c.CopyCacheMap(),
		},
		TraceId: c.TraceId,
	}
}
