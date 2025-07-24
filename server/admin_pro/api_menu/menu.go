package api_menu

import (
	"github.com/77d88/go-kit/plugins/xapi/server/mw"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	pro2 "github.com/77d88/go-kit/server/admin_pro/pro"
	"time"
)

// 所有菜单
func allMenu(c *xhs.Ctx) {
	var menus []*pro2.Menu
	c.Fatalf(xdb.Ctx(c).Find(&menus))
	// 为了返回全部所有菜单权限清空
	c.Send(ConvertMenusToRouter(menus))
}

func saveMenu(c *xhs.Ctx) {
	var req apiRequest
	c.ShouldBind(&req)
	c.Fatalf(req.Name == "", "参数错误")
	c.Fatalf(xdb.Ctx(c).SaveMap(&pro2.Menu{}, req, map[string]interface{}{
		"update_user": c.GetUserIdAssert(),
	}))
}

func deleteMenu(c *xhs.Ctx) {
	var req apiRequest
	c.ShouldBind(&req)
	c.Fatalf(req.Id <= 0, "参数错误")
	c.Fatalf(xdb.Ctx(c).WithId(req.Id).Updates(map[string]interface{}{
		"deleted_time": time.Now(),
		"update_user":  c.GetUserIdAssert(),
	}))
}
func getUserPermissionMenu(c *xhs.Ctx) {
	userId := c.GetUserId()
	var user pro2.User
	if userId <= 0 {
		c.Fatalf(!c.HasRole(pro2.RoleSuperAdmin), "请登录", xhs.FatalWithCode(xhs.CodeTokenError))
	}
	c.Fatalf(xdb.Ctx(c).FindById(&user, userId))
	var menus []*pro2.Menu
	c.Fatalf(xdb.Ctx(c).Find(&menus))
	code := user.AllPermissionCode()
	code = append(code, c.GetRoles()...)
	c.Send(ConvertMenusToRouter(menus, code...))
}

type perRequest struct {
	xhs.IdRequest
	PermissionIds *xdb.Int8Array `json:"permissionIds"`
}

func setPermission(c *xhs.Ctx) {
	var r perRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id == 0, "请选择菜单")

	var menu pro2.Menu
	c.Fatalf(xdb.Ctx(c).FindById(&menu, r.Id))
	c.Fatalf(menu.IsSystem, "系统菜单不能设置")

	var permissions []*pro2.Permission
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Permission{}).WithId(r.PermissionIds.ToSlice()...).Find(&permissions))
	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Menu{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"permission":  xdb.NewTextArray(codes...),
			"update_user": c.GetUserIdAssert(),
		}), "设置权限失败")

}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/all", pro2.SuperAdmin, allMenu)                  // 所有菜单
	api.RegisterPost(path+"/save", pro2.SuperAdmin, saveMenu)                // 保存菜单
	api.RegisterPost(path+"/delete", pro2.SuperAdmin, deleteMenu)            // 删除菜单
	api.RegisterPost(path+"/plist", mw.JwtApiHandler, getUserPermissionMenu) // 获取用户权限菜单
	api.RegisterPost(path+"/smenup", mw.JwtApiHandler, setPermission)        // 设置菜单权限
}

func Default(api *xe.Engine) {
	Register(api, "/x/sys/menu")
}
