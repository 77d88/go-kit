package aliyunossginsts

import (
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts/handler/oss/gettoken"
	"github.com/gin-gonic/gin"
)

func DefaultRegister(path string, x *xe.Engine, handler ...xhs.NewHandlers) {
	x.RegisterByGroup(path, func(r *gin.RouterGroup) {
		r.POST("/getToken", xapi.ApiHandlerToGin(append(handler, gettoken.Run)...)...)
	})
}
