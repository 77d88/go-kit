package aliyunoss

import (
	"context"
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xdb"
	"testing"
)

func Test_Run(t *testing.T) {
	xapi.InitTestConfig()
	xdb.Init()
	Init(nil)
	img, err := genOtherImg(context.TODO(), "temp/666623810818118", 100, 100)
	if err != nil {
		panic(err)
	}
	t.Log(img)

	t.Log(GetOssPostSign(context.TODO()))
}

func Test_OptimizeAll(t *testing.T) {
	xapi.InitTestConfig()
	xdb.Init()
	Init(nil)

	var res []Res
	xdb.Ctx(context.TODO()).Where("id > 0").Find(&res)

	for _, r := range res {
		if r.IsOptimize {
			continue
		}
		err := OptimizeRes(context.TODO(), r, OFile{
			ETag: r.AliEtag,
			Id:   r.ID,
			Key:  r.Path,
			Type: r.MimeType,
		})
		if err != nil {
			return
		}
	}
}
