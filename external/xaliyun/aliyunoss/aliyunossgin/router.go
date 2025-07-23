package aliyunossgin

import (
	ossfilesaveHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/filesave"
	ossgetdomainHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/getdomain"
	ossoptimizeallHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/optimizeall"
	ossPostsign "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/postsign"
	ossPresign "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/presign"
	osssaveHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/save"
	osssavenetlinkHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/savenetlink"
	osssmallimgsaveHandler "github.com/77d88/go-kit/external/xaliyun/aliyunoss/aliyunossgin/handler/oss/smallimgsave"
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/gin-gonic/gin"
)

func DefaultRegister(path string, x *xe.Engine, handler ...xhs.NewHandlers) {
	x.RegisterByGroup(path, func(r *gin.RouterGroup) {
		r.POST("/getDomain", xapi.ApiHandlerToGin(append(handler, ossgetdomainHandler.Run)...)...)       // 获取域名
		r.POST("/save", xapi.ApiHandlerToGin(append(handler, osssaveHandler.Run)...)...)                 // oss保存
		r.POST("/fileSave", xapi.ApiHandlerToGin(append(handler, ossfilesaveHandler.Run)...)...)         // 文件直传
		r.POST("/smallImgSave", xapi.ApiHandlerToGin(append(handler, osssmallimgsaveHandler.Run)...)...) // 小图片文件直传
		r.POST("/saveNetLink", xapi.ApiHandlerToGin(append(handler, osssavenetlinkHandler.Run)...)...)   // 保存网络图片
		r.POST("/optimizeAll", xapi.ApiHandlerToGin(append(handler, ossoptimizeallHandler.Run)...)...)   // 处理所有没有优化的图片 慎用！！！
		r.POST("/postSign", xapi.ApiHandlerToGin(append(handler, ossPostsign.Run)...)...)                // post上传的v4签名
		r.POST("/preSign", xapi.ApiHandlerToGin(append(handler, ossPresign.Run)...)...)                  // put预签名签名上传的地址
	})
}
