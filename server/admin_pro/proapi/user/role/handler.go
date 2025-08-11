package role

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户角色管理
type response struct {
}

type request struct {
	Id    int64          `json:"id,string"`
	Roles *xdb.Int8Array `json:"roles"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Id <= 0 {
		return nil, xerror.New("参数错误:Id不能为空")
	}
	if r.Roles == nil || len(r.Roles.ToSlice()) <= 0 {
		if result := xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
			Updates(map[string]interface{}{
				"permission":       xdb.NewInt8Array(),
				"update_user":      c.GetUserId(),
				"permission_codes": xdb.NewTextArray(),
			}); result.Error != nil {

			return nil, result.Error
		}
		return
	}
	var roles []pro.Role
	if result := xdb.Ctx(c).WithId(r.Roles.ToSlice()...).Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	codes := make([]string, 0, len(roles))
	for _, p := range roles {
		codes = append(codes, p.PermissionCodes.ToSlice()...)
	}

	if result := xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"roles":            r.Roles,
			"update_user":      c.GetUserId(),
			"permission_codes": xdb.NewTextArray(codes...),
		}); result.Error != nil {
		return nil, result.Error
	}

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
