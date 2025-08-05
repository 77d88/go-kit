package create

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/dbmw"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
)

// 创建用户2
type response struct {
}

type request struct {
}

//go:generate xf -m=2
func handler(c *xhs.Ctx, r *request, db *xdb.DB) (resp interface{}, err error) {
	m := make(map[string]interface{})
	// 600075249287237
	scan := db.WithCtx(c).Table("s_user").WithId(60007524928723).Take(&m)
	if scan.IsNotFound() {
		return nil, xerror.New("煤球找到数据")
	}
	return &m, scan.Error
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.GET(path, run(), dbmw.TranManager())
}
