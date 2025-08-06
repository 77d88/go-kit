package menu

import (
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
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
	if result := xdb.Ctx(c).WithId(userId).Take(&user); result.Error != nil {
		return nil, result.Error
	}

	var menus []*pro.Menu

	if user.IsSuperAdmin {
		result := xdb.Ctx(c).Find(&menus)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		result := xdb.Ctx(c).Where("permission in (?)", user.AllPermissionCode()).Find(&menus)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return menus, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
