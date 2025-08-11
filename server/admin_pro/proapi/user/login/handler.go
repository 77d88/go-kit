package login

import (
	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 登录
type response struct {
	*auth.LoginResponse
	Name    string `json:"name"`
	ToSuper bool   `json:"toSuper,omitempty"`
}

type request struct {
	UserName         string `form:"userName" json:"userName"`
	Password         string `form:"password" json:"password"`
	Type             int    `form:"type" json:"type"`
	RememberPassword bool   `form:"rememberPassword" json:"rememberPassword"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	// 获取登录信息
	var manager auth.Manager
	err = c.Server.XE.Invoke(func(ctx auth.Manager) {
		manager = ctx
	})
	if err != nil {
		return nil, xerror.New("获取登录信息失败")
	}
	if r.UserName == "" || r.Password == "" {
		return nil, xerror.New("用户名密码错误")
	}

	switch r.Type {
	case 1:
		return passwordLogin(c, r, manager)
	default:
		return nil, xerror.New("不支持的登录方式")
	}

}

func passwordLogin(c *xhs.Ctx, r *request, manager auth.Manager) (*response, error) {
	pwd := xmd5.EncryptSalt(r.Password, pro.UserPwdSalt)

	// todo 密码改一下 用加密的比较
	if r.UserName == "superadmin" && pwd == "super.admin.(^$@^)@admin.com" {
		authorization, err := manager.Login(-1, auth.WithRoles(pro.Per_SuperAdmin))
		if err != nil {
			return nil, xerror.New("登录失败")
		}
		return &response{
			LoginResponse: authorization,
			ToSuper:       true,
		}, nil
	}

	var user *pro.User
	if result := xdb.Ctx(c).Where(" password = ? and username = ? ", pwd, r.UserName).First(&user); result.Error != nil {
		return nil, xerror.New("用户名或密码错误")
	}
	if user.Disabled {
		xlog.Infof(c, "用户%d 登录失败，用户被禁用", user.ID)
		return nil, xerror.New("请联系管理员")
	}
	login, err := manager.Login(user.ID, auth.WithRoles(user.AllPermissionCode()...))
	if err != nil {
		return nil, xerror.New("登录失败")
	}

	return &response{
		LoginResponse: login,
		Name:          user.Nickname,
	}, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
