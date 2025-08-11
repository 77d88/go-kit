package xjob

import (
	"context"
	"github.com/77d88/go-kit/plugins/xlog"
)

func defaultCtx() context.Context {
	return context.WithValue(context.Background(), xlog.CtxLogParam, map[string]interface{}{
		"origin": "xjob",
	})
}
