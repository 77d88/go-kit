package aliyunaddress

import (
	"context"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
)

func DefaultRegister(path string, x *xhs.HttpServer, handler ...xhs.HandlerMw) {
	init := Init()
	if init == nil {
		xlog.Warnf(context.Background(), "aliyunaddress init error")
		return
	}
	x.POST(path+"/standardizeAddress", Run, handler...) // 地址标准化 默认必须要权限
}
