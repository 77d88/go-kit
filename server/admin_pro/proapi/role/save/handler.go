package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 角色保存
type response struct {
}

type request struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"code,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Name == "" {
		return nil, xerror.New("参数错误:名称不能为空")
	}

	var role pro.Role
	if result := xdb.Ctx(c).Where("name = ?", r.Name).Find(&role); result.Error != nil {
		return nil, result.Error
	}
	if r.Id > 0 && role.ID != r.Id {
		return nil, xerror.New("角色名称不能重复")
	}

	if r.Id > 0 {
		if result := xdb.Ctx(c).Model(&pro.Role{}).WithId(r.Id).Updates(map[string]interface{}{
			"name":        r.Name,
			"update_user": c.GetUserId(),
		}); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if result := xdb.Ctx(c).Model(&pro.Role{}).Create(&pro.Role{
			Name:       r.Name,
			UpdateUser: c.GetUserId(),
		}); result.Error != nil {
			return nil, result.Error
		}
	}
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
