package aliyunaddress

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

func DefaultRegister(path string, x *xhs.HttpServer, handler ...xhs.HandlerMw) {
	x.POST(path+"/standardizeAddress", Run, handler...) // 地址标准化 默认必须要权限
}
