package list

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

// 获取用户列表
type response struct {
}

type request struct {
}

//go:generate xf -m=2
func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.GET(path, run(), auth.ForceAuth)
	xsh.POST(path, run(), auth.ForceAuth)
}
