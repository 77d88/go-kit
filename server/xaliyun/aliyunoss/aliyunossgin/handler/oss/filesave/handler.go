package filesave

import (
	"mime/multipart"

	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// handler 文件直传 /oss/fileSave
func handler(c *xhs.Ctx, r request) (interface{}, error) {
	file, header, err := c.Request.FormFile("file")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	if err != nil {
		return nil, err
	}

	client := aliyunoss.Client
	o := aliyunoss.Config

	key := o.TempPrefix + "\\" + header.Filename

	request := &oss.PutObjectRequest{
		Key:    &key,
		Bucket: &o.OssBucket,
		Body:   file,
	}
	object, err := client.PutObject(c, request)

	if err != nil {
		return nil, err
	}

	etag := object.Headers.Get("ETag")
	return aliyunoss.OFile{ETag: etag, Key: key}, nil

}

// Run 文件直传
func Run(c *xhs.Ctx) (interface{}, error) {
	var r request
	//c.ShouldBind(&r)
	return handler(c, r)
}

type request struct {
}

type response struct {
}
