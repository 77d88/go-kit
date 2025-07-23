package api_user

import (
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/pro"
	"github.com/77d88/go-kit/plugins/xapi/server/mw"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"time"
)

type listRequest struct {
	xdb.ApiPageRequest
	Name     string `json:"name"`
	Disabled *bool  `json:"disabled"`
}

func list(c *xhs.Ctx) {
	var r listRequest
	c.ShouldBind(&r)

	var (
		total int64
		users []*pro.User
	)
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).
		Where("not is_super_admin"). // 超级管理员不显示出来
		XWhere(r.Name != "", "username ilike @name || nickname ilike @name", xdb.Param("name", r.Name)).
		XWhere(r.Disabled != nil, "disabled = ?", r.Disabled).
		IdDesc().FindPage(r, &users, &total))

	c.SendPage(pro.ToUserResponses(users), total)
}

func disable(c *xhs.Ctx) {
	var r xhs.IdRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"disabled":    true,
			"update_user": c.GetUserIdAssert(),
		}), "禁用失败")
}

func del(c *xhs.Ctx) {
	var r xhs.IdRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"update_user":  c.GetUserIdAssert(),
			"deleted_time": time.Now(),
		}), "删除失败")
}
func enable(c *xhs.Ctx) {
	var r xhs.IdRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"disabled":    false,
			"update_user": c.GetUserIdAssert(),
		}), "启用失败")
}

type saveRequest struct {
	Id       int64          `json:"id,string"`
	Password string         `json:"password,omitempty"`
	Disabled bool           `json:"disabled"`
	Username string         `json:"username"`
	Nickname string         `json:"nickname"`
	Avatar   *xdb.Int8Array `json:"avatar"`
	Email    string         `json:"email"`
}

func save(c *xhs.Ctx) {
	var r saveRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Username == "", "请输入用户名")
	if r.Id <= 0 {
		c.Fatalf(r.Password == "", "新用户必须要输入密码")
	}
	c.Fatalf(xdb.Ctx(c).SaveMap(&pro.User{}, r, map[string]interface{}{
		"update_user": c.GetUserIdAssert(),
	}))
}

type permissionRequest struct {
	xhs.IdRequest
	Permission *xdb.Int8Array `json:"permission"`
}

func setPermission(c *xhs.Ctx) {
	var r permissionRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")

	var permissions []pro.Permission
	c.Fatalf(xdb.Ctx(c).Model(&pro.Permission{}).WithId(r.Permission.ToSlice()...).Find(&permissions))
	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}

	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"permission":       r.Permission,
			"update_user":      c.GetUserIdAssert(),
			"permission_codes": xdb.NewTextArray(codes...),
		}), "设置权限失败")
}

type rolesRequest struct {
	xhs.IdRequest
	Roles *xdb.Int8Array `json:"roles"`
}

func setRoles(c *xhs.Ctx) {
	var r rolesRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")

	var codes []string
	if !r.Roles.IsEmpty() { // 获取角色对应的权限码

		var roles []pro.Role
		c.Fatalf(xdb.Ctx(c).Model(&pro.Role{}).WithId(r.Roles.ToSlice()...).Find(&roles))

		var permissionIds []int64
		if len(roles) > 0 {
			for _, role := range roles {
				if role.Permission != nil {
					permissionIds = append(permissionIds, role.Permission.ToSlice()...)
				}
			}
		}

		if len(permissionIds) > 0 {
			var permission []pro.Permission
			c.Fatalf(xdb.Ctx(c).Model(&pro.Permission{}).WithId(permissionIds...).Find(&permission))
			for _, p := range permission {
				codes = append(codes, p.Code)
			}
		}
	}

	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"roles":                 r.Roles,
			"update_user":           c.GetUserIdAssert(),
			"role_permission_codes": xdb.NewTextArray(codes...),
		}),
	)
}

func info(c *xhs.Ctx) {
	var r xhs.IdRequest
	c.ShouldBind(&r)
	c.Fatalf(r.Id <= 0, "请选择用户")
	var user pro.User
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(r.Id).Find(&user))
	c.Send(user.ToResponse())
}

func tokenInfo(c *xhs.Ctx) {
	userId := c.GetUserIdAssert()
	var user pro.User
	c.Fatalf(xdb.Ctx(c).Model(&pro.User{}).WithId(userId).Find(&user))
	c.Send(user.ToResponse())
}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/list", mw.JwtApiHandler, list)
	api.RegisterPost(path+"/save", mw.JwtApiHandler, save)
	api.RegisterPost(path+"/del", mw.JwtApiHandler, del)
	api.RegisterPost(path+"/disable", mw.JwtApiHandler, disable)
	api.RegisterPost(path+"/enable", mw.JwtApiHandler, enable)
	api.RegisterPost(path+"/setRoles", mw.JwtApiHandler, setRoles)
	api.RegisterPost(path+"/suserp", mw.JwtApiHandler, setPermission)
	api.RegisterPost(path+"/info", info)
	api.RegisterPost(path+"/tkInfo", tokenInfo)
	api.RegisterPost(path+"/login", login)
	api.RegisterPost(path+"/refreshToken", refreshToken)
}
