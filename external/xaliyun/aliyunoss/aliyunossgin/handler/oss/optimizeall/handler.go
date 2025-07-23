package smallimgsave

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/external/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
)

// handler 优化所有文件 /oss/optimizeall
func handler(c *xhs.Ctx, r request) {

	var res []aliyunoss.Res
	c.Fatalf(xdb.Ctx(c).Where("id > 0 and mime_type in (0,1) and not is_optimize").Find(&res))

	for _, r := range res {
		if r.IsOptimize {
			continue
		}
		c.Fatalf(aliyunoss.OptimizeRes(c, r, aliyunoss.OFile{
			ETag: r.AliEtag,
			Id:   r.ID,
			Key:  r.Path,
			Type: xcore.Ternary(r.MimeType == 0, aliyunoss.ResTypeImage, r.MimeType),
		}))
	}
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
}
