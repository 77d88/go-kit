package mw

import (
	xhs2 "github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// JwtAuth 鉴权中间件

func JwtApiHandler(c *xhs2.Ctx) {
	if c.GetUserId() == 0 {
		c.SendError(c.NewError("auth error!").SetCode(xhs2.CodeTokenError))
		c.Abort()
		return
	}
	c.Next()
}
