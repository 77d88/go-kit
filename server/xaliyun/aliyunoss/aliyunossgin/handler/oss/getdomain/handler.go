package getdomain

import (
	"fmt"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// handler 获取域名 /oss/getDomain
func handler(c *xhs.Ctx, r request) (interface{}, error) {
	s := aliyunoss.Config
	return &response{
		Domain:       fmt.Sprintf("%s/%s", s.Domain, s.SavePrefix),
		TempPrefix:   s.TempPrefix,
		UploadDomain: s.Domain,
		Region:       s.Region,
	}, nil
}

// Run 获取域名
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	err := c.ShouldBind(&r)
	if err != nil {
		return nil, err
	}
	return handler(c, r)
}

type request struct {
}

type response struct {
	Domain       string `json:"domain" `
	TempPrefix   string `json:"tempPrefix" `
	UploadDomain string `json:"uploadDomain" `
	Region       string `json:"region"`
}
