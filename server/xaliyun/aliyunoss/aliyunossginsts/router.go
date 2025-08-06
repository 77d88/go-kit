package aliyunossginsts

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

func DefaultRegister(path string, x *xhs.HttpServer, handler ...xhs.HandlerMw) {
	x.XE.MustInvoke(func(client *oss.Client) {
		x.POST(path+"/getToken", gettoken.Run, handler...)
	})
}
