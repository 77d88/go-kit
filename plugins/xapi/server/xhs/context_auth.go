package xhs

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
)

type ContextAuth struct {
	userId int64
	roles  []string
	token  string
}

func (c *ContextAuth) SetUserId(id int64) {
	c.userId = id
}

func (c *ContextAuth) GetUserId() int64 {
	return c.userId
}

func (c *ContextAuth) SetRoles(roles ...string) {
	c.roles = roles
}
func (c *ContextAuth) GetRoles() []string {
	return c.roles
}
func (c *ContextAuth) HasRole(role ...string) bool {
	return xarray.ContainAny(c.GetRoles(), role)
}
func (c *ContextAuth) HasRolesAll(roles ...string) bool {
	return xarray.ContainAll(c.GetRoles(), roles)
}

func (c *ContextAuth) SetToken(token string) {
	c.token = token
}
func (c *ContextAuth) GetToken() string {
	return c.token
}

func (c *Ctx) GetUserIdAssert(msg ...string) int64 {
	userId := c.GetUserId()
	if userId > 0 {
		return userId
	}
	if len(msg) == 0 {
		panic(xerror.New("用户异常").SetCode(CodeTokenError))
	} else {
		panic(xerror.New(msg[0]).SetCode(CodeTokenError))
	}
}
