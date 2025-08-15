package login

import (
	"testing"

	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
)

func TestName(t *testing.T) {

	println(xmd5.Encrypt("super.admin.(^$@^)@admin.com"))

}
