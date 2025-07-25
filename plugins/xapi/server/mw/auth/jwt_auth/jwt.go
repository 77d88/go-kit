package jwt_auth

import (
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var (
	// JwtKey jwtKey 随意即可
	jwtKey        = []byte(xmd5.Encrypt("jwtKey@$%!23"))
	refreshJwtKey = []byte(xmd5.Encrypt("refreshJwtKey@$%!23"))
)

type JwtAuth struct {
	key         []byte
	refreshKey  []byte
	AutoRenewal bool
}

func New() *JwtAuth {
	return &JwtAuth{
		key:         jwtKey,
		refreshKey:  refreshJwtKey,
		AutoRenewal: true,
	}
}
func NewCustomize(key, refreshKey []byte) *JwtAuth {
	return &JwtAuth{
		key:         key,
		refreshKey:  refreshKey,
		AutoRenewal: true,
	}
}

func (a *JwtAuth) GenerateToken(id int64, expr time.Duration, roles ...string) (string, error) {
	return generateToken(id, expr, a.key, roles...)
}
func (a *JwtAuth) GenerateRefreshToken(id int64, expr time.Duration, roles ...string) (string, error) {
	return generateToken(id, expr, a.refreshKey, roles...)
}
func (a *JwtAuth) VerificationToken(jwtStr string) *auth.VerificationData {
	return verificationToken(jwtStr, a.key)
}
func (a *JwtAuth) VerificationRefreshToken(token string) *auth.VerificationData {
	return verificationToken(token, a.refreshKey)
}
func (a *JwtAuth) SetAutoRenewal(autoRenewal bool) *JwtAuth {
	a.AutoRenewal = autoRenewal
	return a
}

func (a *JwtAuth) IsAutoRenewal() bool {
	return a.AutoRenewal
}

// Login api登录
func (a *JwtAuth) Login(id int64, roles ...string) (*auth.LoginResponse, error) {
	// 生成一个短期有效的token 10分钟
	token, err := a.GenerateToken(id, time.Minute*10, roles...)
	if err != nil {
		return nil, err
	}
	// 生成一个长期有效的token 30 天
	longToken, err := a.GenerateToken(id, time.Hour*24*30, roles...)
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponse{
		Token:        token,
		RefreshToken: longToken,
	}, nil
}

// verificationToken 校验token
func verificationToken(jwtStr string, key []byte) *auth.VerificationData {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return &auth.VerificationData{
			Err: xerror.Newf("jwt parse error %s", err),
		}
	}
	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return &auth.VerificationData{
			Err: xerror.Newf("jwt Claims error %s", err),
		}
	}
	if !token.Valid {
		return &auth.VerificationData{
			ExpireTime: expirationTime.Time,
			Err:        xerror.New("jwt Valid error "),
		}
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return &auth.VerificationData{
			ExpireTime: expirationTime.Time,
			Err:        xerror.Newf("jwt GetSubject error %v", err),
		}
	}
	content, i, err := parseContent(subject)
	return &auth.VerificationData{
		Id:         content,
		Roles:      i,
		ExpireTime: expirationTime.Time,
	}
}

// generateToken 生成token
// xid 用户id
// roles 角色
// expr 过期时间
func generateToken(id int64, expr time.Duration, key []byte, roles ...string) (string, error) {
	exp := time.Now().Add(expr).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "jx",
		"sub": buildContent(id, roles),
		"exp": exp,
	})
	return token.SignedString(key)
}

func buildContent(id int64, roles []string) string {
	return fmt.Sprintf("%d&%s", id, xarray.Join(roles, ","))
}

func parseContent(str string) (int64, []string, error) {
	split := strings.Split(str, "&")
	id, err := xstr.ToInt[int64](split[0])
	if err != nil {
		return 0, nil, err
	}
	return id, xstr.SplitAndTrim(split[1], ",", " "), nil
}
