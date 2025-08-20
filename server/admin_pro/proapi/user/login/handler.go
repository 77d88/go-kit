package login

import (
	"github.com/77d88/go-kit/basic/xencrypt/xpwd"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
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
	manager, err := x.Get[auth.Manager]()
	if err != nil {
		return nil, xerror.New("获取登录信息失败")
	}
	if r.UserName == "" || r.Password == "" {
		return nil, xerror.New("用户名密码错误")
	}

	if r.Type == 0 {
		r.Type = 1
	}

	switch r.Type {
	case 1:
		return passwordLogin(c, r, manager)
	default:
		return nil, xerror.New("不支持的登录方式")
	}

}

func passwordLogin(c *xhs.Ctx, r *request, manager auth.Manager) (*response, error) {
	pwd, err := xpwd.HashPassword(r.Password)
	if err != nil {
		return nil, xerror.New("密码错误")
	}
	if r.UserName == "superadmin" && xpwd.CheckPasswordHash(r.Password, "$2a$10$vkeRXagMQVyizbROxJMkE.2WTUvsp8E.pqPBUiqpkeszUfwvEtMMq") {
		authorization, err := manager.Login(c, -1, auth.WithRoles(pro.Per_SuperAdmin), auth.WithSinglePoint())
		if err != nil {
			xlog.Debugf(c, "登录失败 %v", err)
			return nil, xerror.New("登录失败")
		}
		return &response{
			LoginResponse: authorization,
			ToSuper:       true,
		}, nil
	}

	var user *pro.User
	if result := xdb.C(c).Where(" password = ? and username = ? ", pwd, r.UserName).First(&user); result.Error != nil {
		return nil, xerror.New("用户名或密码错误")
	}
	if user.Disabled {
		xlog.Infof(c, "用户%d 登录失败，用户被禁用", user.ID)
		return nil, xerror.New("请联系管理员")
	}
	login, err := manager.Login(c, user.ID, auth.WithRoles(user.AllPermissionCode()...))
	if err != nil {
		return nil, xerror.New("登录失败")
	}

	return &response{
		LoginResponse: login,
		Name:          user.Nickname,
	}, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run())
}
