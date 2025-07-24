package xhs

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xsys"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
)

type Handler func(ctx *Ctx) (interface{}, error)

// WarpHandle 通用处理函数 包装了一个本地的Context
func WarpHandle(f Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := f(newCtx(c))
		if err != nil {
			handleError(newCtx(c), err)
		}
	}
}

func handleError(c *Ctx, e interface{}) {
	switch r := e.(type) {
	case string, int:
		xlog.Errorf(c, "%v", r)
		c.SendError(xerror.New(r))
	case Response, *Response: // 理论上不会走到这里 这里只有错误
		c.Send(r)
	case xerror.Error, *xerror.Error:
		c.SendError(r)
	case error:
		// 正常业务的错误不会走到这里 打印错误
		xlog.Errorf(c, "%v %s", r, xsys.StackTrace(false))
		c.SendError(xerror.New(r))
	default:
		c.SendError(xerror.New(r))
	}
}

// NewHandlers api执行处理器 包括异常 事务
func NewHandlers(c *Ctx, fs ...Handler) {
	defer func() {
		if !c.Writer.Written() {
			// 如果没有写入内容，则默认返回成功
			if c.Result == nil {
				c.JSON(CodeSuccess, NewResp(nil))
				return
			}
			c.JSON(CodeSuccess, c.Result)
		}
	}()
	defer func() {
		if e := recover(); e != nil {
			handleError(c, e)
		}
	}()
	// 处理业务逻辑
	for _, f := range fs {
		r, err := f(c)
		if err != nil {
			handleError(c, err)
			break
		}
		if r != nil {
			c.Result = NewResp(r)
		}
	}

	// 继续处理
	c.Next()

}
