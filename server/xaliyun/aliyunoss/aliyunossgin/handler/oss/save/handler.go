package save

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// Run oss保存
func Run(c *xhs.Ctx) {
	var r file
	c.ShouldBind(&r)
	c.Fatalf(r.Id == 0 && r.Key == "", "文件列表为空")
	handler(c, r)
}

// handler oss保存 /oss/save
func handler(c *xhs.Ctx, r file) {
	save, err := aliyunoss.DbSave(c, r.OFile)
	c.Fatalf(err)
	c.Send(save)
}

type file struct {
	aliyunoss.OFile
}
