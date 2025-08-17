package list

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 角色列表
type response struct {
}

type request struct {
	Page xdb.PageSearch `json:"page"`
	Name string         `json:"name"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	db := xdb.XWhere(xdb.C(c), r.Name != "", "name ilike @name", xdb.Param("name", xdb.WarpLike(r.Name)))
	if result := xdb.FindPage[pro.Role](db, r.Page, true); result.Error != nil {
		return nil, result.Error
	} else {
		return xhs.NewResp(result.List, result.Total), nil
	}
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
