package list

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 权限列表
type response struct {
	Id   int64  `json:"id,string"`
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type request struct {
	xdb.PageSearch
	Code string `json:"code,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var permission []pro.Permission
	var total int64
	if result := xdb.Ctx(c).XWhere(r.Code != "", "code ilike @code", xdb.Param("code", xdb.WarpLike(r.Code))).FindPage(r, &permission, &total); result.Error != nil {
		return nil, result.Error
	}
	return xhs.NewResp(xarray.Map(permission, func(index int, item pro.Permission) response {
		return response{
			Id:   item.ID,
			Code: item.Code,
			Desc: item.Desc,
		}
	}), total), nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
