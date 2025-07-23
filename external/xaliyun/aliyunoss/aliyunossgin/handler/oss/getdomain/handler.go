package getdomain

import (
	"fmt"
	"github.com/77d88/go-kit/external/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// handler 获取域名 /oss/getDomain
func handler(c *xhs.Ctx, r request) {
	s := aliyunoss.Config
	c.Send(&response{
		Domain:       fmt.Sprintf("%s/%s", s.Domain, s.SavePrefix),
		TempPrefix:   s.TempPrefix,
		UploadDomain: s.Domain,
		Region:       s.Region,
	})
}

// Run 获取域名
func Run(c *xhs.Ctx) {
	var r request
	c.ShouldBind(&r)
	handler(c, r)
}

type request struct {
}

type response struct {
	Domain       string `json:"domain" `
	TempPrefix   string `json:"tempPrefix" `
	UploadDomain string `json:"uploadDomain" `
	Region       string `json:"region"`
}
