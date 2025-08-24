package list

import (
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 角色列表
type response struct {
}

type request struct {
	Page xdb.PageSearch   `json:"page"`
	Name string           `json:"name"`
	Ids  xtype.Int64Array `json:"ids"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var roles []pro.Role
	result := xpg.C(c).Model(&pro.Role{}).XWhere(r.Name != "", "name ilike ?", xdb.WarpLike(r.Name)).
		XWhere(len(r.Ids) > 0, "id = any(?)", r.Ids).
		FindPage(&roles, r.Page, len(r.Ids) == 0)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return xhs.NewResp(roles, result.Total), nil
	}
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/role/list", run(), auth.ForceAuth)
}
