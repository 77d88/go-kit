package aliyunaddress

import (
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xapi/server/mw"
	handler2 "github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/gin-gonic/gin"
)

func DefaultRegister(path string, x *xe.Engine, handler ...handler2.NewHandlers) {
	x.RegisterByGroup(path, func(r *gin.RouterGroup) {
		r.POST("/standardizeAddress", xapi.ApiHandlerToGin(append(handler, mw.JwtApiHandler, Run)...)...) // 地址标准化 默认必须要权限
	})
}
