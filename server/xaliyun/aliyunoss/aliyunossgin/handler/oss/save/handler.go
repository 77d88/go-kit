package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
)

// Run oss保存
func Run(c *xhs.Ctx) (interface{}, error) {
	var r file
	err := c.ShouldBind(&r)
	if err != nil {
		return nil, err
	}
	if r.Id == 0 && r.Key == "" {
		return nil, xerror.New("文件列表为空")
	}
	return handler(c, r)
}

// handler oss保存 /oss/save
func handler(c *xhs.Ctx, r file) (interface{}, error) {
	return aliyunoss.DbSave(c, xpg.C(c), r.OFile)
}

type file struct {
	aliyunoss.OFile
}
