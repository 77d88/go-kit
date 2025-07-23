package filesave

import (
	aliyunoss2 "github.com/77d88/go-kit/external/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// handler 文件直传 /oss/fileSave
func handler(c *xhs.Ctx, r request) {
	file, header, err := c.Request.FormFile("file")
	defer file.Close()
	c.Fatalf(err)
	client := aliyunoss2.Client
	o := aliyunoss2.Config
	c.Fatalf(err)
	key := o.TempPrefix + "\\" + header.Filename

	request := &oss.PutObjectRequest{
		Key:    &key,
		Bucket: &o.OssBucket,
		Body:   file,
	}
	object, err := client.PutObject(c, request)
	c.Fatalf(err)
	etag := object.Headers.Get("ETag")
	c.Send(aliyunoss2.OFile{ETag: etag, Key: key})

}

// Run 文件直传
func Run(c *xhs.Ctx) {
	var r request
	//c.ShouldBind(&r)
	handler(c, r)
}

type request struct {
}

type response struct {
}
