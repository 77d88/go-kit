package xmd5

import (
	"crypto/md5"
	"fmt"
	"github.com/77d88/go-kit/basic/xstr"
	"strings"
)

// Encrypt 计算字符串MD5
func Encrypt(str string) string {
	bytes := []byte(str)
	digest := md5.New()
	digest.Write(bytes)
	sum := digest.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

// EncryptSalt 使用md5+salt 加密密码
func EncryptSalt(pwd string, salt string) string {
	if xstr.IsBlank(pwd) {
		return ""
	}
	return strings.ToLower(Encrypt(pwd + salt))
}
