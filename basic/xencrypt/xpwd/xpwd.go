package xpwd

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // 可调整 cost
	return string(bytes), err
}

// CheckPasswordHash 对比明文密码与哈希值
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
