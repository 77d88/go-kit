package aliyunossginsts

import (
	"context"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func DefaultRegister(path string, s *xhs.HttpServer, handler ...xhs.HandlerMw) {
	err := x.Find(func(client *oss.Client) {
		if client == nil {
			xlog.Warnf(context.TODO(), "oss sts client is nil oss route notRegister")
			return
		}
		s.POST(path+"/getToken", gettoken.Run, handler...)
	})
	if err != nil {
		panic(err)
	}
}
