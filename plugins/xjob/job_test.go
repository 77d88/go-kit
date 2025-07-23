package xjob

import (
	"context"
	"github.com/77d88/go-kit/plugins/xlog"
	"testing"
)

func TestName2(t *testing.T) {
	Init()
	Submit(context.WithValue(context.Background(), xlog.CtxLogParam, map[string]interface{}{
		"origin": "xjobasasfasf",
	}), func(c context.Context) {
		xlog.Errorf(c, "test")
	})
	select {}
}

func TestName(t *testing.T) {
	SubmitCron("test", "*/1 * * * * *", func() {
		xlog.Infof(nil, "test")
	})
	select {}

}
