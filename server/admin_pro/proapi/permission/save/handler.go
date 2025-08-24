package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 权限保存
type response struct {
}

type request struct {
	Id   int64  `json:"id,omitempty,string"`
	Code string `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Code == "" {
		return nil, xerror.New("权限码不能为空")
	}
	if r.Code == pro.Per_SuperAdmin {
		return nil, xerror.New("该权限码已被占用")
	}
	var permission pro.Permission
	if result := xpg.C(c).Where("code = ?", r.Code).Find(&permission); result.Error != nil {
		return nil, result.Error
	}
	if permission.ID > 0 && permission.ID != r.Id {
		return nil, xerror.New("权限码已存在")
	}
	var old pro.Permission
	if r.Id > 0 {
		if result := xpg.C(c).Where("id = ?", r.Id).First(&old); result.Error != nil {
			return nil, result.Error
		}
	}

	err = xpg.C(c).Transaction(func(tx *xpg.Inst) error {
		// 修改权限本身
		if result := tx.Model(&pro.Permission{}).Save(r, func(m map[string]interface{}) {
			m["update_user"] = c.GetUserId()
		}); result.Error != nil {
			return result.Error
		}

		if old.ID > 0 && r.Code != old.Code { // 修改权限码的情况
			args := map[string]interface{}{
				"code":    old.Code,
				"newCode": r.Code,
				"codeArr": xdb.NewTextArray(old.Code),
			}
			// 更新角色里面的权限码
			if result := tx.Exec(`update s_sys_role set 
                      "permission_codes" = array_append(array_remove("permission_codes", @code), @newCode )
                  where deleted_time is null and permission_codes && @codeArr `, args); result.Error != nil {
				return result.Error
			}
			// 更新用户里面的权限码
			if result := tx.Exec(`update s_sys_user set 
                      "permission_codes" = array_append(array_remove("permission_codes", @code), @newCode )
                  where deleted_time is null  and permission_codes && @codeArr  `, args); result.Error != nil {
				return result.Error
			}
			// 更新用户里面的权限码2
			if result := tx.Exec(`update s_sys_user set 
                      "role_permission_codes" = array_append(array_remove("role_permission_codes", @code), @newCode )
                  where deleted_time is null  and role_permission_codes && @codeArr  `, args); result.Error != nil {
				return result.Error
			}
		}
		return nil
	})

	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/permission/save", run(), auth.ForceAuth)
}
