package xpwd

import (
	"testing"

	"github.com/77d88/go-kit/basic/xtime"
)

func TestName(t *testing.T) {
	inv := xtime.NewTimeInterval()
	t.Log(inv.IntervalMs())
	password := Password("123456")
	t.Log(inv.IntervalMs())
	t.Log(password)
	t.Log(CheckPassword("123456", password))
	t.Log(inv.IntervalMs())
}
