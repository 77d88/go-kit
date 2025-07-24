package smallimgsave

import (
	"context"
	"encoding/base64"
	"fmt"
	context2 "github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"mime/multipart"
)

// handler 小图片文件直传 /oss/smallImgSave
func handler(c *context2.Ctx, r request) {
	file, header, err := c.Request.FormFile("file")
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.SendError(err)
		}
	}(file)
	c.Fatalf(err)
	client := aliyunoss.Client
	o := aliyunoss.Config
	c.Fatalf(err)
	key := o.TempPrefix + "/" + header.Filename

	request := &oss.PutObjectRequest{
		Bucket: &o.OssBucket,
		Key:    &key,
		Body:   file,
	}
	object, err := client.PutObject(c, request)
	c.Fatalf(err)
	etag := object.Headers.Get("ETag")

	save, err := aliyunoss.DbSave(c, aliyunoss.OFile{
		ETag: etag + "small", // 小图片的etag加上small标识 不要影响原始图片
		Key:  key,
	}, func(c context.Context, key, toPath string) error {
		targetImageName := toPath
		style := "image/resize,m_fixed,w_200,h_200" // 小图片都处理为200x200
		process := fmt.Sprintf("%s|sys/saveas,o_%v,b_%v", style,
			base64.URLEncoding.EncodeToString([]byte(targetImageName)),
			base64.URLEncoding.EncodeToString([]byte(o.OssBucket)))

		request := &oss.ProcessObjectRequest{
			Bucket:  oss.Ptr(o.OssBucket), // 指定要操作的存储空间名称
			Key:     oss.Ptr(key),         // 指定要处理的图片名称
			Process: oss.Ptr(process),     // 指定处理指令
		}
		_, err2 := client.ProcessObject(c, request)
		return err2
	})
	c.Fatalf(err)
	c.Send(save)
}

// Run 小图片文件直传
func Run(c *context2.Ctx) {
	var r request
	//c.ShouldBind(&r)
	handler(c, r)
}

type request struct {
}

type response struct {
}
