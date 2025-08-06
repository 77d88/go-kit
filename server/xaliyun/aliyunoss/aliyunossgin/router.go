package aliyunossgin

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	ossfilesaveHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/filesave"
	ossgetdomainHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/getdomain"
	ossoptimizeallHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/optimizeall"
	ossPostsign "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/postsign"
	ossPresign "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/presign"
	osssaveHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/save"
	osssavenetlinkHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/savenetlink"
	osssmallimgsaveHandler "github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin/handler/oss/smallimgsave"
)

func DefaultRegister(path string, r *xhs.HttpServer, handler ...xhs.HandlerMw) {

	//x.RegisterByGroup(path, func(r *gin.RouterGroup) {
	r.POST("/getDomain", ossgetdomainHandler.Run, handler...)       // 获取域名
	r.POST("/save", osssaveHandler.Run, handler...)                 // oss保存
	r.POST("/fileSave", ossfilesaveHandler.Run, handler...)         // 文件直传
	r.POST("/smallImgSave", osssmallimgsaveHandler.Run, handler...) // 小图片文件直传
	r.POST("/saveNetLink", osssavenetlinkHandler.Run, handler...)   // 保存网络图片
	r.POST("/optimizeAll", ossoptimizeallHandler.Run, handler...)   // 处理所有没有优化的图片 慎用！！！
	r.POST("/postSign", ossPostsign.Run, handler...)                // post上传的v4签名
	r.POST("/preSign", ossPresign.Run, handler...)                  // put预签名签名上传的地址
	//})
}
