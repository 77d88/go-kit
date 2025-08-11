package aliyunossginsts

import (
	"context"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
)

func DefaultRegister(path string, s *xhs.HttpServer, handler ...xhs.HandlerMw) {
	client := aliyunoss.InitWith()
	if client == nil {
		xlog.Warnf(context.TODO(), "oss client is nil oss route notRegister")
		return
	}
	s.POST(path+"/getToken", gettoken.Run, handler...)
}
