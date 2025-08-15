package aliyunossginsts

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
)

func DefaultRegister(path string, s *xhs.HttpServer, handler ...xhs.HandlerMw) {
	s.POST(path+"/getToken", gettoken.Run, handler...)
}
