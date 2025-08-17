package del

import (
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
	"gorm.io/gorm"
)

// 权限删除
type response struct {
}

type request struct {
	Id   int64  `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Id <= 0 {
		return nil, xerror.New("参数错误:Id不能为空")
	}
	var permission pro.Permission
	if result := xdb.C(c).Where("id = ?", r.Id).First(&permission); result.Error != nil {
		return nil, result.Error
	}

	if err := xdb.C(c).Transaction(func(d *gorm.DB) error {
		// 删除权限
		if result := xdb.Session(d).Model(&pro.Permission{}).Where("id = ?", r.Id).
			Updates(map[string]interface{}{
				"update_user":  c.GetUserId(),
				"deleted_time": time.Now(),
			}); result.Error != nil {
			return result.Error
		}
		// 移除角色里面的权限
		if result := xdb.Session(d).Exec(`update s_sys_role set 
                      "permission" = array_remove("permission", @id) ,
                      "permission_codes" = array_remove("permission_codes", @code) 
                  where deleted_time is null `, xdb.Param("id", r.Id), xdb.Param("code", permission.Code)); result.Error != nil {
			return result.Error
		}

		// 移除用户里面的权限码
		if result := xdb.Session(d).Exec(`update s_sys_user set 
                      "permission_codes" = array_remove("permission_codes", @code) ,
                      "role_permission_codes" = array_remove("role_permission_codes", @code) 
                  where deleted_time is null `, xdb.Param("code", permission.Code)); result.Error != nil {
			return result.Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
