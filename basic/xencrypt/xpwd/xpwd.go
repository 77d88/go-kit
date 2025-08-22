package xpwd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/77d88/go-kit/basic/xsys"
	"github.com/77d88/go-kit/plugins/xlog"
)

// SecretKey 密钥 不要暴露给外部
var SecretKey = []byte("default-secret-key-change-in-production")

func init() {
	// 从环境变量获取密钥
	key := xsys.OsEnvGet("SECRET_KEY", "default-secret-key-change-in-production")
	SecretKey = []byte(key)
	if key == "default-secret-key-change-in-production" {
		xlog.Errorf(nil, "secret key is not set, using default key")
	}
}

// Password 对密码进行哈希 如果有更高安全性要求 "golang.org/x/crypto/bcrypt" 性能较低
func Password(password string) string {
	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}

// CheckPassword 对比明文密码与哈希值
func CheckPassword(password, hash string) bool {
	expectedHash := Password(password)
	return hmac.Equal([]byte(hash), []byte(expectedHash))
}
