package create

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// Create
type response struct {
}

type request struct {
}

//go:generate xf -m=2
func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run())
}
