package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
)

const (
	HttpHeaderKey = "Authorization"
)

type Option struct {
	Roles       []string      `json:"roles,omitempty"`
	Data        interface{}   `json:"data,omitempty"`
	Duration    time.Duration `json:"-"`
	MaxLoginNum uint          `json:"-"` // 最大登录数量
}

type OptionHandler interface {
	Apply(manager *Option)
}

type OptionFunc func(manager *Option)

func (f OptionFunc) Apply(manager *Option) {
	f(manager)
}
func WithRoles(roles ...string) OptionHandler {
	return OptionFunc(func(manager *Option) {
		manager.Roles = roles
	})
}
func WithData(data interface{}) OptionHandler {
	return OptionFunc(func(manager *Option) {
		manager.Data = data
	})
}
func WithDuration(duration time.Duration) OptionHandler {
	return OptionFunc(func(manager *Option) {
		manager.Duration = duration
	})
}
func WithMaxLoginNum(maxLoginNum uint) OptionHandler {
	return OptionFunc(func(manager *Option) {
		if maxLoginNum == 0 {
			maxLoginNum = 1
		}
		// 默认这里设置不能超过100
		if maxLoginNum > 100 {
			maxLoginNum = 100
		}
		manager.MaxLoginNum = maxLoginNum
	})
}

func GetOpt(handler ...OptionHandler) *Option {
	opt := Option{
		Roles:    []string{},
		Duration: time.Hour * 24 * 7,
		Data:     nil,
	}
	for _, f := range handler {
		if option, ok := f.(OptionHandler); ok {
			option.Apply(&opt)
		}
	}
	return &opt
}

type Manager interface {
	GenerateToken(id int64, opt ...OptionHandler) (string, error)
	GenerateRefreshToken(id int64, opt ...OptionHandler) (string, error)
	VerificationToken(token string) *VerificationData
	VerificationRefreshToken(token string) *VerificationData
	Login(id int64, opt ...OptionHandler) (*LoginResponse, error)
	Logout(token string) error
	IsAutoRenewal() bool
}

type ApiAuth struct {
	Manager Manager
}

func NewX() *ApiAuth {
	get, err := x.Get[Manager]()
	if err != nil {
		xlog.Errorf(nil, "get auth manager error: %v", err)
		return nil
	}
	return &ApiAuth{
		Manager: get,
	}
}

func New(manager Manager) *ApiAuth {
	return &ApiAuth{
		Manager: manager,
	}
}

func NewMw(manager Manager) xhs.HandlerMw {
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
	Data       interface{} // 附带 数据
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
func (c *ApiAuth) GenerateToken(id int64, opt ...OptionHandler) (string, error) {
	return c.Manager.GenerateToken(id, opt...)
}

// VerificationToken 验证token获取信息
func (c *ApiAuth) VerificationToken(token string) *VerificationData {
	return c.Manager.VerificationToken(token)
}

// Authorization 通用授权
func (c *ApiAuth) Authorization(id int64, opt ...OptionHandler) (*LoginResponse, error) {
	return c.Manager.Login(id, opt...)
}

// TokenInfo 默认的token解析中间件 只负责鉴权获取用户信息 不负责强制验证登录授权
func (c *ApiAuth) TokenInfo() xhs.HandlerMw {
	return func(x *xhs.Ctx) {
		if c.Manager == nil {
			x.Next()
			return
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
			x.Next()
			return
		}
		data := manager.VerificationToken(token)
		if !data.Validate() {
			x.Next()
			return
		}
		x.SetUserId(data.Id)
		x.SetPermission(data.Roles...)
		x.SetToken(token)

		x.Next()
	}
}

func ForceAuth(c *xhs.Ctx) {
	if c.GetUserId() == 0 {
		c.Send(xerror.New("登录异常!").SetCode(xhs.CodeTokenError))
		c.Abort()
		return
	}
	c.Next()
}
