package loginInfo

import (
	"time"

	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// 用户登录信息
type response struct {
	UserId     int64     `json:"userId,omitempty"`
	Roles      []string  `json:"roles,omitempty"`
	ExpireTime time.Time `json:"expireTime"`
}

func handler(c *xhs.Ctx) (interface{}, error) {
	return &response{
		UserId:     c.Auth.UserId,
		Roles:      c.Auth.Roles,
		ExpireTime: c.Auth.ExpireTime,
	}, nil
}

//func run() xhs.Handler {
//	return func(ctx *xhs.Ctx) (interface{}, error) {
//		return handler(ctx, nil)
//	}
//}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/loginInfo", handler, auth.ForceAuth)
}
