package list

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 菜单列表
type response struct {
}

type request struct {
	ParentId int64 `json:"parentId,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.ParentId == 0 {
		return xpg.C(c).Where("root_menu").Order("sort asc,id asc").Find(&[]pro.Menu{}).Decon()
	}

	var menu pro.Menu
	if result := xpg.C(c).Where("id = ?", r.ParentId).First(&menu); result.Error != nil {
		if result.IsNotFound() {
			return make([]struct{}, 0), nil
		} else {
			return nil, result.Error
		}
	}
	if len(menu.Children) == 0 {
		return make([]struct{}, 0), nil
	}
	var menus []pro.Menu
	if result := xpg.C(c).Where("id = ANY(?)", menu.Children).Order("sort asc,id asc").Find(&menus); result.Error != nil {
		if result.IsNotFound() {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return menus, nil
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/menu/list", run(), auth.ForceAuth, pro.HansPermission(pro.Per_SuperAdmin))
}
