package save

import (
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// Run oss保存
func Run(c *xhs.Ctx) {
	var r file
	c.ShouldBind(&r)
	c.Fatalf(len(r.Files) == 0, "文件列表为空")
	handler(c, r)
}

// handler oss保存 /oss/save
func handler(c *xhs.Ctx, r file) {
	var mm map[string]string
	for _, v := range r.Files {
		save, err := aliyunoss.DbSave(c, v)
		c.Fatalf(err)
		mm[v.ETag] = xstr.ToString(save)
	}
	c.Send(mm)
}

type file struct {
	Files []aliyunoss.OFile `form:"files" json:"files"`
}
