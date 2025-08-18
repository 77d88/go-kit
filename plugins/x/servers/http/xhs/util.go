package xhs

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xerror"
)

func ParamZeroCheck(check ...interface{}) error {
	for i, v := range check {
		if xcore.IsZero(v) {
			return xerror.Newf("参数错误200%d", i).SetCode(CodeParamError)
		}
	}
	return nil
}
