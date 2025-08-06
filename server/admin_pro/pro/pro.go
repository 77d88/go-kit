package pro

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// Per_SuperAdmin 特殊权限超级管理员
const Per_SuperAdmin = "superAdmin"

// HansPermissionAny 权限处理 满足任意一个
func HansPermissionAny(permission ...string) xhs.HandlerMw {
	return func(c *xhs.Ctx) {
		if !c.HasPermission(permission...) {
			c.SendError(xerror.New("权限不足").SetCode(xhs.CodeNoPermissions))
			c.Abort()
			return
		}

	}
}

// HansPermission 权限处理 满足所有
func HansPermission(permission ...string) xhs.HandlerMw {
	return func(c *xhs.Ctx) {
		if !c.HasPermission(permission...) {
			c.SendError(xerror.New("权限不足").SetCode(xhs.CodeNoPermissions))
			c.Abort()
			return
		}

	}
}