package menu

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户菜单列表
type request struct {
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	userId := c.GetUserId()
	var menus []*pro.Menu
	if c.Auth.HasPermission(pro.Per_SuperAdmin) {
		result := xpg.C(c).Find(&menus)
		if result.Error != nil {
			return nil, result.Error
		} else {
			return menus, nil
		}
	}

	var user pro.User
	if result := xpg.C(c).Where("id = ?", userId).First(&user); result.Error != nil {
		return nil, result.Error
	}

	result := xpg.C(c).Where("permission in (?)", user.AllPermissionCode()).Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/menu/menu", run(), auth.ForceAuth)
}
