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
	Page xdb.PageSearch `json:"page"`
	Code string         `json:"code,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	db := xdb.C(c)
	if r.Code != "" {
		db = db.Where("code ilike @code", xdb.Param("code", xdb.WarpLike(r.Code)))
	}
	if result := xdb.FindPage[pro.Permission](db, r.Page, true); result.Error != nil {
		return nil, result.Error
	} else {
		return xhs.NewResp(xarray.Map(result.List, func(index int, item pro.Permission) response {
			return response{
				Id:   item.ID,
				Code: item.Code,
				Desc: item.Desc,
			}
		}), result.Total), nil
	}

}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/permission/list", run(), auth.ForceAuth)
}
