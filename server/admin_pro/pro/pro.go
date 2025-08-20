package pro

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// Per_SuperAdmin 特殊权限超级管理员
const Per_SuperAdmin = "superAdmin"

func HansPermission(permission ...string) xhs.HandlerMw {
	return func(c *xhs.Ctx) {
		if !c.Auth.HasPermission(permission...) {
			c.SendError(xerror.New("权限不足").SetCode(xhs.CodeNoPermissions))
			c.Abort()
			return
		}

	}
}
