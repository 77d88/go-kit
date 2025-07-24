package api_role

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/xapi/server/mw"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	pro2 "github.com/77d88/go-kit/server/admin_pro/pro"
	"time"
)

type request struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"code,omitempty"`
}

func save(c *xhs.Ctx) {
	var r request
	c.ShouldBind(&r)
	c.Fatalf(r.Name == "", "请输入权限码")

	var role pro2.Role
	c.Fatalf(xdb.Ctx(c).Where("name = ?", r.Name).Find(&role))
	if role.ID > 0 {
		c.Fatalf(role.ID != r.Id, "角色名称不能重复")
	}

	c.Fatalf(xdb.Ctx(c).Model(&pro2.Role{}).WithId(r.Id).Updates(map[string]interface{}{
		"name":        r.Name,
		"update_user": c.GetUserIdAssert(),
	}))
}

func del(c *xhs.Ctx) {
	var r request
	c.ShouldBind(&r)
	c.Fatalf(r.Id == 0, "参数错误")
	var role pro2.Permission
	c.Fatalf(xdb.Ctx(c).WithId(r.Id).First(&role))
	var count int64
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Role{}).Where("deleted_time is null and roles && ?", xdb.NewInt8Array(r.Id)).Count(&count))
	c.Fatalf(count > 0, "角色已分配给用户，不能删除")
	// 删除角色
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Role{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"update_user":  c.GetUserIdAssert(),
			"deleted_time": time.Now(),
		}))
}

type listRequest struct {
	xdb.ApiPageRequest
	Name string `form:"name"`
}

func list(c *xhs.Ctx) {
	var r listRequest
	c.ShouldBind(&r)
	var (
		list  []pro2.Role
		total int64
	)
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Permission{}).Where("deleted_time is null").
		XWhere(r.Name != "", "code ilike @name", xdb.Param("code", xdb.WarpLike(r.Name))).FindPage(r, &list, &total))

	c.SendPage(xarray.Map(list, func(index int, item pro2.Role) *pro2.RoleDst {
		return item.ToResponse()
	}), total)
}

type permissionRequest struct {
	xhs.IdRequest
	Permission *xdb.Int8Array `json:"permission"`
}

func setPermission(c *xhs.Ctx) {
	var r permissionRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择角色")

	var permissions []pro2.Permission
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Permission{}).WithId(r.Permission.ToSlice()...).Find(&permissions))
	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}

	c.Fatalf(xdb.Ctx(c).Model(&pro2.Role{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"permission":       r.Permission,
			"update_user":      c.GetUserIdAssert(),
			"permission_codes": xdb.NewTextArray(codes...),
		}), "设置权限失败")
}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/list", mw.JwtApiHandler, list)
	api.RegisterPost(path+"/save", mw.JwtApiHandler, save)
	api.RegisterPost(path+"/del", mw.JwtApiHandler, del)
	api.RegisterPost(path+"/srolep", mw.JwtApiHandler, del)
}
