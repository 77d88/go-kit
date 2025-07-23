package aes_auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"io"
	"strings"
	"time"
)

var (
	defaultAuthV2Key        = []byte("asdzxcvbnm,./..-defaultAuthV2Key")
	defaultRefreshAuthV2Key = []byte("asd][zva-defaultRefreshAuthV2Key")
)

type AesAuth struct {
	key         []byte // 资源密钥
	refreshKey  []byte // 刷新token的秘钥
	AutoRenewal bool
}

func (a *AesAuth) GenerateToken(id int64, expr time.Duration, roles ...string) (string, error) {
	return generateToken(id, expr, a.key, roles...)
}
func (a *AesAuth) GenerateRefreshToken(id int64, expr time.Duration, roles ...string) (string, error) {
	return generateToken(id, expr, a.refreshKey, roles...)
}

func (a *AesAuth) VerificationToken(jwtStr string) *auth.VerificationData {
	return verificationToken(jwtStr, a.key)
}
func (a *AesAuth) VerificationRefreshToken(token string) *auth.VerificationData {
	return verificationToken(token, a.refreshKey)
}

// Login api登录
func (a *AesAuth) Login(id int64, roles ...string) (*auth.LoginResponse, error) {
	// 生成一个短期有效的token 10分钟
	token, err := a.GenerateToken(id, time.Minute*30, roles...)
	if err != nil {
		return nil, err
	}
	// 生成一个长期有效的token 30 天
	longToken, err := a.GenerateToken(id, time.Hour*24*30, roles...)
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponse{
		Id:           id,
		Token:        token,
		RefreshToken: longToken,
	}, nil
}

func (a *AesAuth) SetAutoRenewal(autoRenewal bool) *AesAuth {
	a.AutoRenewal = autoRenewal
	return a
}
func (a *AesAuth) IsAutoRenewal() bool {
	return a.AutoRenewal
}

func New() *AesAuth {
	return &AesAuth{
		key:        defaultAuthV2Key,
		refreshKey: defaultRefreshAuthV2Key,
	}
}
func NewCustomize(key, refreshKey []byte) *AesAuth {
	return &AesAuth{
		key:        key,
		refreshKey: refreshKey,
	}
}

// GenerateToken 生成token 第二版本 使用AES加密 暴露时间戳 用户端可自行校验是否过期 base64 无填充请注意
func generateToken(id int64, expr time.Duration, key []byte, roles ...string) (string, error) {
	exp := time.Now().Add(expr).UTC().Unix()
	content := buildContent(id, roles)
	// content 加密
	// 1. 创建AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 2. 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. 生成随机Nonce（GCM要求Nonce唯一）
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 4. 加密数据（Seal方法返回nonce + ciphertext）
	ciphertext := gcm.Seal(nonce, nonce, []byte(content), nil)

	// 5. 返回Base64编码的Token
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(fmt.Sprintf("%s|%d", ciphertext, exp))), nil
}

// DecryptTokenV2 解密Base64编码的Token，返回原始数据
func verificationToken(token string, key []byte) *auth.VerificationData {

	// 1. Base64解码
	ciphertext, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
	if err != nil {
		return &auth.VerificationData{
			Err: err,
		}
	}
	// 2. 拆分Token
	parts := strings.Split(string(ciphertext), "|")

	ciphertext = []byte(parts[0])
	if len(parts) != 2 {
		return &auth.VerificationData{
			Err: xerror.New("invalid token length error"),
		}
	}

	exp, err := xparse.ToNumber[int64](parts[1])

	if err != nil {
		return &auth.VerificationData{
			Err: xerror.New("invalid token time error"),
		}
	}

	// 3. 创建AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return &auth.VerificationData{
			ExpireTime: time.Unix(exp, 0),
			Err:        xerror.New("invalid token create cipher error"),
		}
	}

	// 4. 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return &auth.VerificationData{
			ExpireTime: time.Unix(exp, 0),
			Err:        xerror.New("invalid token NewGCM error"),
		}
	}

	// 5. 提取Nonce（GCM要求Nonce在密文前）
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return &auth.VerificationData{
			ExpireTime: time.Unix(exp, 0),
			Err:        xerror.New("invalid token GCM Nonce error"),
		}
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 6. 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return &auth.VerificationData{
			ExpireTime: time.Unix(exp, 0),
			Err:        xerror.New("invalid token GCM Open error"),
		}
	}

	// 7. 返回id和role
	id, roles, err := parseContent(string(plaintext))

	black := auth.IsUserBlack(id)
	if black { // 用户被拉黑了 不支持针对单token的禁用
		return &auth.VerificationData{
			Id:         id,
			Roles:      roles,
			ExpireTime: time.Unix(exp, 0),
			Err:        xerror.New("请联系管理员"),
		}
	}

	return &auth.VerificationData{
		Id:         id,
		Roles:      roles,
		ExpireTime: time.Unix(exp, 0),
	}
}

func buildContent(id int64, roles []string) string {
	return fmt.Sprintf("%d|%s", id, xarray.Join(roles, ","))
}

func parseContent(str string) (int64, []string, error) {
	split := xstr.Split(str, "|")
	id, err := xparse.ToNumber[int64](split[0])
	if err != nil {
		return 0, nil, err
	}
	return id, xstr.SplitAndTrim(split[1], ",", " "), nil
}
