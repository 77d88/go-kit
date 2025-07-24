package pro

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/bits-and-blooms/bloom/v3"
)

// 构建能够接收 10 万个元素且误报率为 1% 的 Bloom 过滤器。
var filter = bloom.NewWithEstimates(10000, 0.01)

const RoleSuperAdmin = "superAdmin"

func SuperAdmin(c *xhs.Ctx) {
	if c.GetUserId() == 0 {
		c.SendError(c.NewError("auth error!").SetCode(xhs.CodeTokenError))
		c.Abort()
		return
	}

	if !c.HasRole(RoleSuperAdmin) { // 只有超级管理员 权限里面才会有这个
		c.SendError(c.NewError("auth error by role!").SetCode(xhs.CodeNoPermissions))
		c.Abort()
		return
	}

	c.Next()
}
