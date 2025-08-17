package list

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 菜单列表
type response struct {
}

type request struct {
	ParentId int64 `json:"parentId,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var menu pro.Menu
	if result := xdb.C(c).Where("id = ?", r.ParentId).Take(&menu); result.Error != nil {
		if xdb.IsNotFound(result.Error) {
			return make([]struct{}, 0), nil
		} else {
			return nil, result.Error
		}
	}
	if menu.Children.IsEmpty() {
		return make([]struct{}, 0), nil
	}
	var menus []pro.Menu
	if result := xdb.C(c).Where("id in ?", menu.Children.ToSlice()).Find(&menus); result.Error != nil {
		if xdb.IsNotFound(result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return menus, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
