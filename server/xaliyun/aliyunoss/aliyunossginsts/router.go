package aliyunossginsts

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
)

func DefaultRegister(path string, x *xhs.HttpServer, handler ...xhs.HandlerMw) {
		x.POST(path+"/getToken",gettoken.Run, handler...)
}
