package smallimgsave

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	aliyunoss2 "github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// handler 优化所有文件 /oss/optimizeall
func handler(c *xhs.Ctx, r request) (interface{}, error) {

	var res []aliyunoss2.Res
	if result := xdb.Ctx(c).Where("id > 0 and mime_type in (0,1) and not is_optimize").Find(&res); result.Error != nil {

		return nil, result.Error
	}

	for _, r := range res {
		if r.IsOptimize {
			continue
		}
		err := aliyunoss2.OptimizeRes(c, r, aliyunoss2.OFile{
			ETag: r.AliEtag,
			Id:   r.ID,
			Key:  r.Path,
			Type: xcore.Ternary(r.MimeType == 0, aliyunoss2.ResTypeImage, r.MimeType),
		})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
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
}
