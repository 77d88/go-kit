package auth

import (
	"errors"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
	"net/http"
	"time"
)

const (
	HttpHeaderKey = "Authorization"
)

type Manager interface {
	GenerateToken(id int64, expr time.Duration, roles ...string) (string, error)
	GenerateRefreshToken(id int64, expr time.Duration, roles ...string) (string, error)
	VerificationToken(jwtStr string) *VerificationData
	VerificationRefreshToken(token string) *VerificationData
	Login(id int64, roles ...string) (*LoginResponse, error)
	IsAutoRenewal() bool
}

type ApiAuth struct {
	Manager Manager
}

func New(manager Manager) *ApiAuth {
	return &ApiAuth{
		Manager: manager,
	}
}

func NewMw(manager Manager) xhs.Handler {
	auth := ApiAuth{
		Manager: manager,
	}
	return auth.TokenInfo()
}

type LoginResponse struct {
	Id           int64  `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type VerificationData struct {
	Id         int64     `json:"id"`
	Roles      []string  `json:"roles"`
	ExpireTime time.Time `json:"expireTime"`
	Err        error
}

// Validate  验证token
func (d *VerificationData) Validate() bool {

	// 有异常
	if d.Err != nil {
		return false
	}

	// 无效的ID
	if d.Id == 0 {
		return false
	}

	// 时间验证
	if d.IsExpired() {
		return false
	}
	return true
}

// IsExpired 是否过期
func (d *VerificationData) IsExpired() bool {
	if d.ExpireTime.IsZero() || time.Now().After(d.ExpireTime) {
		return true
	}
	return false

}

// GenerateToken 生成token
func (c *ApiAuth) GenerateToken(id int64, expr time.Duration, roles ...string) (string, error) {
	return c.Manager.GenerateToken(id, expr, roles...)
}

// VerificationToken 验证token获取信息
func (c *ApiAuth) VerificationToken(token string) *VerificationData {
	return c.Manager.VerificationToken(token)
}

// Authorization 通用授权
func (c *ApiAuth) Authorization(id int64, roles ...string) (*LoginResponse, error) {
	return c.Manager.Login(id, roles...)
}

// TokenInfo 默认的token解析中间件 只负责鉴权获取用户信息 不负责强制验证登录授权
func (c *ApiAuth) TokenInfo() xhs.Handler {
	return func(x *xhs.Ctx) (interface{}, error) {
		if c.Manager == nil {
			return nil, nil
		}
		manager := c.Manager
		// 常规校验
		token := x.Query(HttpHeaderKey)
		if token == "" {
			token = x.GetHeader(HttpHeaderKey) // 再从头部获取一下
		}
		if token == "" {
			cookie, err := x.Cookie(HttpHeaderKey)
			if err != nil && !errors.Is(err, http.ErrNoCookie) {
				xlog.Errorf(x, "获取cookie失败 %s", err)
			} else {
				if cookie != "" {
					token = cookie
				}
			}
		}
		if token == "" {
			return nil, nil
		}
		data := manager.VerificationToken(token)
		if !data.Validate() {
			return nil, nil
		}
		x.SetUserId(data.Id)
		x.SetRoles(data.Roles...)
		x.SetToken(token)
		return nil, nil
	}
}

func ForceAuth(c *xhs.Ctx) (interface{}, error) {
	if c.GetUserId() == 0 {
		return nil, c.NewError("auth error!").SetCode(xhs.CodeTokenError)
	}
	return nil, nil
}
