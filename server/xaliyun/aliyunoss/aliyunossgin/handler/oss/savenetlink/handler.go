package savenetlink

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
)

// handler 保存网络图片 /oss/saveNetLink
func handler(c *xhs.Ctx, r request) (interface{}, error) {

	//client := ossutilt.Client
	//o := ossutilt.Config
	//bucket, err := client.Bucket(o.OssBucket)
	//c.Fatalf(err)
	//
	//key := o.TempPrefix + "/" + uuid.NewString()
	//
	//options := make([]oss.Option, 0)
	//opts := oss.AddContentType(options, key)
	//
	//// 下载网络图片
	//
	////file, err := oss.OpenUrl(client, r.Url)
	//
	//request := &oss.PutObjectRequest{
	//	ObjectKey: key,
	//	Reader:    file,
	//}
	//object, err := bucket.DoPutObject(request, opts)
	//c.Fatalf(err)
	//etag := object.Headers.Get("ETag")
	return nil, nil
}

// Run 保存网络图片
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	return handler(c, r)
}

type request struct {
	Url string `form:"url" json:"url"`
}

type response struct {
}
