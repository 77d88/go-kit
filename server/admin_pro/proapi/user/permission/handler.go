package permission

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户授权管理
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
	if r.Permission == nil || len(r.Permission.ToSlice()) <= 0 {
		if result := xdb.C(c).Model(&pro.User{}).Where("id = ?", r.Id).
			Updates(map[string]interface{}{
				"permission":       xdb.NewInt8Array(),
				"update_user":      c.GetUserId(),
				"permission_codes": xdb.NewTextArray(),
			}); result.Error != nil {

			return nil, result.Error
		}
		return
	}
	codes := make([]string, 0)
	if !r.Permission.IsEmpty() {
		var permission []pro.Permission
		if result := xdb.C(c).Where("id in ?", r.Permission.ToSlice()).Find(&permission); result.Error != nil {
			return nil, result.Error
		}
		for _, p := range permission {
			codes = append(codes, p.Code)
		}
	}

	if result := xdb.C(c).Model(&pro.User{}).Where("id = ?", r.Id).
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
