package xhs

import (
	"github.com/77d88/go-kit/basic/xarray"
)

type ContextAuth struct {
	userId     int64
	permission []string
	token      string
}

func (c *ContextAuth) SetUserId(id int64) {
	c.userId = id
}

func (c *ContextAuth) GetUserId() int64 {
	return c.userId
}

func (c *ContextAuth) SetPermission(permission ...string) {
	c.permission = permission
}
func (c *ContextAuth) GetPermission() []string {
	return c.permission
}
func (c *ContextAuth) HasPermission(permission ...string) bool {
	return xarray.ContainAny(c.GetPermission(), permission)
}
func (c *ContextAuth) HasPermissionAll(permission ...string) bool {
	return xarray.ContainAll(c.GetPermission(), permission)
}

func (c *ContextAuth) SetToken(token string) {
	c.token = token
}
func (c *ContextAuth) GetToken() string {
	return c.token
}
