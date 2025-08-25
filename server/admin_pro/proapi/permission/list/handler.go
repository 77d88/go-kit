package list

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 权限列表
type response struct {
	Id   int64  `json:"id,string"`
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type request struct {
	Page xpg.PageSearch   `json:"page"`
	Code string           `json:"code,omitempty"`
	Ids  xtype.Int64Array `json:"ids"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var pres []pro.Permission
	result := xpg.C(c).Model(&pro.Permission{}).
		XWhere(r.Code != "", "code ilike ?", xpg.WarpLike(r.Code)).
		XWhere(len(r.Ids) > 0, "id =  any(?)", r.Ids).
		FindPage(&pres, r.Page, len(r.Ids) == 0)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return xhs.NewResp(xarray.Map(pres, func(index int, item pro.Permission) response {
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
