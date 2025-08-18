package create

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

type request struct {
}

func Handler(ctx *xhs.Ctx, r *request) (resp interface{}, err error) {
	return
}

// Create
func run() xhs.Handler {
	return xhs.DefaultShouldHandler[request](Handler)
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run())
}
