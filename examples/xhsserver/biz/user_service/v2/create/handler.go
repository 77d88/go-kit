package create

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

// 创建用户2
type response struct {
}

type request struct {
}

//go:generate xf -m=2
func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	m := make(map[string]interface{})
	// 600075249287237
	scan := xdb.C(c).Table("s_user").Where("id = 60007524928723").Take(&m)
	if xdb.IsNotFound(scan.Error) {
		return nil, xerror.New("煤球找到数据")
	}
	return &m, scan.Error
}

func Register(xsh *xhs.HttpServer) {
	xsh.GET(path, run())
}
