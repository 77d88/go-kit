package save

import (
	"github.com/77d88/go-kit/basic/xencrypt/xpwd"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户保存
type response struct {
}

type request struct {
	Id       int64            `json:"id,string"`
	Password string           `json:"password,omitempty"`
	Disabled bool             `json:"disabled"`
	Username string           `json:"username"`
	Nickname string           `json:"nickname"`
	Avatar   xtype.Int64Array `json:"avatar"`
	Email    string           `json:"email"`
	Roles    xtype.Int64Array `json:"roles"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Username == "" {
		return nil, xerror.New("用户名不能为空")
	}

	if r.Id <= 0 && r.Password == "" {
		return nil, xerror.New("新用户必须要输入密码")
	}

	if r.Password != "nil" {
		password := xpwd.Password(r.Password)
		r.Password = password
	}

	var user pro.User
	if result := xpg.C(c).Where("username = ?", r.Username).Find(&user); result.Error != nil {
		return nil, result.Error
	}
	if user.ID > 0 && user.ID != r.Id {
		return nil, xerror.New("用户名已存在")
	}
	var role []pro.Role
	var roleIds []int64
	var roleNames []string
	var roleCodes []string
	if !r.Roles.IsEmpty() {
		if result := xpg.C(c).Where("id = any(?)", r.Roles).Find(&role); result.Error != nil {
			return nil, result.Error
		}
		for _, p := range role {
			roleIds = append(roleIds, p.ID)
			roleCodes = append(roleCodes, p.PermissionCodes...)
			roleNames = append(roleNames, p.Name)
		}
	}

	if result := xpg.C(c).Model(&user).Save(r, func(m map[string]interface{}) {
		m["update_user"] = c.GetUserId()
		m["roles"] = roleIds
		m["role_names"] = roleNames
		m["role_permission_codes"] = roleCodes
	}); result.Error != nil {
		return nil, result.Error
	}
	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/save", run(), auth.ForceAuth)
}
