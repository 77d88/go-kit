package api_user

import (
	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xapi/ctx"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	pro2 "github.com/77d88/go-kit/server/admin_pro/pro"
	"time"
)

type loginRequest struct {
	UserName         string `form:"userName" json:"userName"`
	Password         string `form:"password" json:"password"`
	Type             int    `form:"type" json:"type"`
	RememberPassword bool   `form:"rememberPassword" json:"rememberPassword"`
}

type loginResponse struct {
	*auth.LoginResponse
	Name    string `json:"name"`
	ToSuper bool   `json:"toSuper,omitempty"`
}

func login(c *xhs.Ctx) {

	var r loginRequest
	c.ShouldBind(&r)

	if r.Type == 1 && r.UserName == "superadmin" && r.Password == "superadmin@admin.com" {
		authorization, err := ctx.Authorization(-1, pro2.RoleSuperAdmin)
		c.Fatalf(err)
		c.Send(&loginResponse{
			LoginResponse: authorization,
			ToSuper:       true,
		})
		return
	}

	pwd := xmd5.EncryptSalt(r.Password, pro2.UserPwdSalt)

	var user pro2.User
	result := xdb.Ctx(c).Where(" password = ? and username = ? ", pwd, r.UserName).First(&user)
	c.Fatalf(result, "用户名或密码错误")
	c.Fatalf(user.Disabled, "请联系管理员", "用户 %d 已被禁用 尝试登录", user.ID)
	code := user.AllPermissionCode()
	if user.IsSuperAdmin {
		code = append(code, pro2.RoleSuperAdmin)
	}

	res, err := ctx.Authorization(user.ID, code...)
	c.Fatalf(err, "登录失败")
	c.Send(loginResponse{
		LoginResponse: res,
		Name:          user.Nickname,
	})
}

func refreshToken(c *xhs.Ctx) {
	request := struct {
		Token    string `json:"token"`
		Duration uint64 `json:"duration"` // 时长 单位秒
	}{}
	c.ShouldBind(&request)
	c.Fatalf(request.Token == "", "获取授权失败")
	if request.Duration <= 0 {
		request.Duration = 60 * 30 // 默认30分钟
	}
	c.Fatalf(request.Duration > 60*60*24, "The duration cannot exceed 1 day")
	data := xapi.Server.AuthManager.VerificationToken(request.Token)
	if !data.Validate() {
		c.SendError(xerror.New("error").SetCode(xhs.CodeRefreshTokenError))
		return
	}

	// 刷新token
	generateToken, err := xapi.Server.AuthManager.GenerateToken(data.Id, time.Second*time.Duration(request.Duration), data.Roles...)
	c.Fatalf(err)
	token := generateToken
	refreshToken := request.Token

	if xapi.Server.AuthManager.IsAutoRenewal() {
		// 如果刷新token的时效小于2小时，则刷新token
		if data.ExpireTime.Sub(time.Now()) < time.Hour*2 {
			// 刷新token
			generateToken, err := xapi.Server.AuthManager.GenerateRefreshToken(data.Id, time.Second*time.Duration(request.Duration), data.Roles...)
			c.Fatalf(err)
			refreshToken = generateToken
		}
	}

	c.Send(&auth.LoginResponse{
		Id:           data.Id,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
