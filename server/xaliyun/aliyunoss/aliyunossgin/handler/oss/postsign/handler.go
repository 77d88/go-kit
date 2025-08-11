package smallimgsave

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// handler 预签名url /oss/postsign 不可以正常使用put上传的用post签名上传 比如微信小程序
func handler(c *xhs.Ctx, r request) (interface{}, error) {
	return aliyunoss.GetOssPostSign(c)
}

// Run 小图片文件直传
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	//c.ShouldBind(&r)
	return handler(c, r)
}

type request struct {
}

type response struct {
	Url string `json:"url"`
}
