package smallimgsave

import (
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"time"
)

// handler 预签名url /oss/presign 可以正常使用put上传的都行
func handler(c *xhs.Ctx, r request) (interface{}, error) {
	key := oss.Ptr(aliyunoss.Config.TempPrefix + "/" + xid.NextIdStr())
	client := aliyunoss.Client
	// 生成PutObject的预签名URL
	result, err := client.Presign(c, &oss.PutObjectRequest{
		Bucket:      oss.Ptr(aliyunoss.Config.OssBucket),
		Key:         key,
		ContentType: r.ContentType,
	},
		oss.PresignExpires(10*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"url":           result.URL,
		"key":           key,
		"signedHeaders": result.SignedHeaders,
	}, err
}

// Run 小图片文件直传
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	err := c.ShouldBind(&r)
	if err != nil {
		return nil, err
	}
	return handler(c, r)
}

type request struct {
	ContentType *string `json:"contentType"`
}

type response struct {
}
