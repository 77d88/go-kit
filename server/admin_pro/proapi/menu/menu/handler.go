package menu

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户菜单列表
type response struct {
}

type request struct {
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	userId := c.GetUserId()

	var user pro.User
	if result := xdb.C(c).Where("id = ?", userId).Take(&user); result.Error != nil {
		return nil, result.Error
	}

	var menus []*pro.Menu

	if user.IsSuperAdmin {
		result := xdb.C(c).Find(&menus)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := xdb.C(c).Where("permission in (?)", user.AllPermissionCode()).Find(&menus)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return menus, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
