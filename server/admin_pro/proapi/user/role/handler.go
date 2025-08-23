package role

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户角色管理
type response struct {
}

type request struct {
	Id    int64            `json:"id,string"`
	Roles xtype.Int64Array `json:"roles"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Id <= 0 {
		return nil, xerror.New("参数错误:Id不能为空")
	}
	if r.Roles == nil || len(r.Roles.ToSlice()) <= 0 {
		if result := xpg.C(c).Model(&pro.User{}).Where("id = ?", r.Id).
			Updates(map[string]interface{}{
				"permission":       nil,
				"update_user":      c.GetUserId(),
				"permission_codes": nil,
			}); result.Error != nil {

			return nil, result.Error
		}
		return
	}
	var roles []pro.Role
	if result := xpg.C(c).Where("id = any(?)", r.Roles).Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	codes := make([]string, 0, len(roles))
	for _, p := range roles {
		codes = append(codes, p.PermissionCodes...)
	}

	if result := xpg.C(c).Model(&pro.User{}).Where("id = ?", r.Id).
		Updates(map[string]interface{}{
			"roles":            r.Roles,
			"update_user":      c.GetUserId(),
			"permission_codes": xdb.NewTextArray(codes...),
		}); result.Error != nil {
		return nil, result.Error
	}

	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/role", run(), auth.ForceAuth)
}
