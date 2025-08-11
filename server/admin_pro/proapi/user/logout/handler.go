package logout

import (
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// 登出/踢出
type response struct {
}

type request struct {
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	// 获取登录信息
	var manager auth.Manager
	err = x.Find(func(ctx auth.Manager) {
		manager = ctx
	})
	return nil, manager.Logout(c.GetToken())
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
