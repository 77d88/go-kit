package xhs

import (
	"time"

	"github.com/77d88/go-kit/basic/xarray"
)

func (c *Ctx) GetUserId() int64 {
	if c.Auth == nil {
		return 0
	}
	return c.Auth.UserId
}

type ContextAuth struct {
	UserId     int64       `json:"userId,omitempty"`
	Roles      []string    `json:"roles,omitempty"`
	Token      string      `json:"token,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	ExpireTime time.Time   `json:"expireTime"`
}

func (c *ContextAuth) HasPermission(permission ...string) bool {
	return xarray.ContainAny(c.Roles, permission)
}
func (c *ContextAuth) HasPermissionAll(permission ...string) bool {
	return xarray.ContainAll(c.Roles, permission)
}
