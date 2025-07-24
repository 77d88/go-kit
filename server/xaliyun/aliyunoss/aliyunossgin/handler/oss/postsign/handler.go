package smallimgsave

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// handler 预签名url /oss/postsign 不可以正常使用put上传的用post签名上传 比如微信小程序
func handler(c *xhs.Ctx, r request) {
	body, err := aliyunoss.GetOssPostSign(c)
	c.Fatalf(err)
	c.Send(body)
}

// Run 小图片文件直传
func Run(c *xhs.Ctx) {
	var r request
	//c.ShouldBind(&r)
	handler(c, r)
}

type request struct {
}

type response struct {
	Url string `json:"url"`
}
