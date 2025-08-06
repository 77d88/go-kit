package del

import (
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// 权限删除
type response struct {
}

type request struct {
}

//go:generate xf -m=2
func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, auth.ForceAuth, run())
}
