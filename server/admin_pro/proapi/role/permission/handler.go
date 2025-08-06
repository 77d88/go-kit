package permission

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 角色授权管理
type response struct {
}

type request struct {
	Id         int64          `json:"id,string"`
	Permission *xdb.Int8Array `json:"permission"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Id <= 0 {
		return nil, xerror.New("参数错误:Id不能为空")
	}
	var permissions []pro.Permission
	if result := xdb.Ctx(c).FindByIds(&permissions, r.Permission.ToSlice()...); result.Error != nil {
		return nil, result.Error
	}
	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}
	codes = xarray.Union(codes)
	if result := xdb.Ctx(c).Model(&pro.Role{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"permission":       r.Permission,
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
