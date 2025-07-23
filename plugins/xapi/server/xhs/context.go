package xhs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
)

const ctxKeyContext = "CTX_KEY_CONTEXT"

type Ctx struct {
	*gin.Context
	ContextAuth
	ApiCache
	Result     interface{}
	PrintStack bool
	TraceId    int64
	test       bool
}

type CopyContext struct {
	context.Context
	ContextAuth
	ApiCache
	TraceId int64
}

func newCtx(c *gin.Context) *Ctx {
	value := c.Value(ctxKeyContext)
	if value != nil {
		c2, ok := value.(*Ctx)
		if ok {
			return c2
		}
	}
	d := &Ctx{
		Context:    c,
		PrintStack: true,
		Result:     nil,
		TraceId:    xid.NextId(),
	}
	d.Set(ctxKeyContext, d)
	return d
}

func NewTestContext() *Ctx {
	return &Ctx{
		test:    true,
		Context: &gin.Context{},
	}
}

func (c *Ctx) Value(key any) any {
	switch key {
	case xlog.CtxLogParam:
		return logFields(c)
	case ctxKeyContext:
		return c
	}
	return c.Context.Value(key)
}

func (c *Ctx) Send(v interface{}) {
	c.Result = &Response{
		Code: CodeSuccess,
		Msg:  "ok",
		Data: v,
	}
	if c.test {
		indent, err := json.MarshalIndent(c.Result, "", "  ")
		if err != nil {
			fmt.Printf("test result:\n%v\n", c.Result)
		} else {
			fmt.Printf("test result:\n%v\n", string(indent))
		}
	}
}

func (c *Ctx) FastSend(v interface{}) {
	c.Fatalf(&Response{
		Code: CodeSuccess,
		Msg:  "ok",
		Data: v,
	})
}

func (c *Ctx) SendJSON(obj any) {
	c.Result = obj
}

func (c *Ctx) SendPage(result interface{}, total int64) {
	c.Result = &Response{
		Code:  CodeSuccess,
		Msg:   "ok",
		Data:  result,
		Total: xcore.V2p(int(total)),
	}
}

func (c *Ctx) FastSendPage(result interface{}, total int64) {
	c.Fatalf(&Response{
		Code:  CodeSuccess,
		Msg:   "ok",
		Data:  result,
		Total: xcore.V2p(int(total)),
	})
}

// ShouldBind 重写一下
func (c *Ctx) ShouldBind(obj any) {
	if obj == nil {
		return
	}
	err := c.Context.ShouldBind(obj)
	if err != nil {
		xlog.Errorf(c, "参数错误 ShouldBind: %+v", err)
		c.Fatalf(xerror.Newf("参数错误").SetCode(CodeParamError))
	}
}

// CopyFinal 拷贝最终上下文 用于传输获其他现场调用
func (c *Ctx) CopyFinal() context.Context {
	return CopyContext{
		Context:     context.WithValue(context.Background(), xlog.CtxLogParam, logFields(c)),
		ContextAuth: c.ContextAuth,
		ApiCache: ApiCache{
			cache: c.CopyCacheMap(),
		},
		TraceId: c.TraceId,
	}
}
