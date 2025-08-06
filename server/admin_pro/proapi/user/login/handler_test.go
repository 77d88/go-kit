package login

import (
	"github.com/77d88/go-kit/basic/xencrypt/xmd5"
	"testing"
)

func TestName(t *testing.T) {


	println(xmd5.Encrypt("super.admin.(^$@^)@admin.com"))

}