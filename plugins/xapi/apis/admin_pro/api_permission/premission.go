package api_permission

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/pro"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"time"
)

type request struct {
	Id   int64  `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}

func save(c *xhs.Ctx) {
	var r request
	c.ShouldBind(&r)
	c.Fatalf(r.Code == "", "请输入权限码")

	var permission pro.Permission
	c.Fatalf(xdb.Ctx(c).Where("code = ?", r.Code).Find(&permission))
	if permission.ID > 0 {
		c.Fatalf(permission.ID != r.Id, "权限码已存在")
	}
	var old pro.Permission
	if r.Id > 0 {
		c.Fatalf(xdb.Ctx(c).WithId(r.Id).First(&old))
	}

	err := xdb.CtxTran(c, func(d *xdb.DataSource) error {
		// 修改权限本身
		if result := d.Session().SaveMap(&pro.Permission{}, r, map[string]interface{}{
			"update_user": c.GetUserIdAssert(),
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
			if result := d.Session().Exec(`update s_sys_role set 
                      "permission_codes" = array_append(array_remove("permission_codes", @code), @newCode )
                  where deleted_time is null and permission_codes && @codeArr `, args); result.Error != nil {
				return result.Error
			}
			// 更新用户里面的权限码
			if result := d.Session().Exec(`update s_sys_user set 
                      "permission_codes" = array_append(array_remove("permission_codes", @code), @newCode )
                  where deleted_time is null  and permission_codes && @codeArr  `, args); result.Error != nil {
				return result.Error
			}
			// 更新用户里面的权限码2
			if result := d.Session().Exec(`update s_sys_user set 
                      "role_permission_codes" = array_append(array_remove("role_permission_codes", @code), @newCode )
                  where deleted_time is null  and role_permission_codes && @codeArr  `, args); result.Error != nil {
				return result.Error
			}
		}

		return nil

	})

	c.Fatalf(err)
}

func del(c *xhs.Ctx) {
	var r request
	c.ShouldBind(&r)
	c.Fatalf(r.Id == 0, "参数错误")
	var permission pro.Permission
	c.Fatalf(xdb.Ctx(c).WithId(r.Id).First(&permission))

	c.Fatalf(xdb.CtxTran(c, func(d *xdb.DataSource) error {
		// 删除权限
		if result := d.Session().Model(&pro.Permission{}).WithId(r.Id).
			Updates(map[string]interface{}{
				"update_user":  c.GetUserIdAssert(),
				"deleted_time": time.Now(),
			}); result.Error != nil {
			return result.Error
		}
		// 移除角色里面的权限
		if result := d.Session().Exec(`update s_sys_role set 
                      "permission" = array_remove("permission", @id) ,
                      "permission_codes" = array_remove("permission_codes", @code) 
                  where deleted_time is null `, xdb.Param("id", r.Id), xdb.Param("code", permission.Code)); result.Error != nil {
			return result.Error
		}

		// 移除用户里面的权限码
		if result := d.Session().Exec(`update s_sys_user set 
                      "permission_codes" = array_remove("permission_codes", @code) ,
                      "role_permission_codes" = array_remove("role_permission_codes", @code) 
                  where deleted_time is null `, xdb.Param("code", permission.Code)); result.Error != nil {
			return result.Error
		}
		return nil
	}))
}

type listRequest struct {
	xdb.ApiPageRequest
	Code string `form:"code"`
}

func list(c *xhs.Ctx) {
	var r listRequest
	c.ShouldBind(&r)
	var (
		list  []pro.Permission
		total int64
	)
	c.Fatalf(xdb.Ctx(c).Model(&pro.Permission{}).Where("deleted_time is null").
		XWhere(r.Code != "", "code ilike @code", xdb.Param("code", xdb.WarpLike(r.Code))).FindPage(r, &list, &total))

	c.SendPage(xarray.Map(list, func(index int, item pro.Permission) *pro.PermissionDst {
		return item.ToResponse()
	}), total)
}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/list", pro.SuperAdmin, list)
	api.RegisterPost(path+"/save", pro.SuperAdmin, save)
	api.RegisterPost(path+"/del", pro.SuperAdmin, del)
}
