package del

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// 用户删除
type response struct {
}

type request struct {
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/del", run(), auth.ForceAuth)
}
