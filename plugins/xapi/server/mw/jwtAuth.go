package mw

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// JwtAuth 鉴权中间件

func JwtApiHandler(c *xhs.Ctx) {
	if c.GetUserId() == 0 {
		c.SendError(c.NewError("auth error!").SetCode(xhs.CodeTokenError))
		c.Abort()
		return
	}
	c.Next()
}
