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
	xdb.PageSearch
	Name string `json:"name"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var (
		list  []pro.Role
		total int64
	)
	if result := xdb.Ctx(c).Model(&pro.Role{}).
		XWhere(r.Name != "", "name ilike ?", xdb.WarpLike(r.Name)).FindPage(r, &list, &total); result.Error != nil {

		return nil, result.Error
	}
	return xhs.NewResp(list, total), nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
