package list

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典列表

type request struct {
	ParentId int64 `json:"parentId,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.ParentId == 0 {
		return xdb.G[pro.Dict]().Where("root").Order("sort asc,id asc").Find(c)
	}

	var dict pro.Dict
	if result := xdb.C(c).Where("id = ?", r.ParentId).Take(&dict); result.Error != nil {
		if xdb.IsNotFound(result.Error) {
			return make([]struct{}, 0), nil
		} else {
			return nil, result.Error
		}
	}
	if dict.Children.IsEmpty() {
		return make([]struct{}, 0), nil
	}
	var dicts []pro.Dict
	if result := xdb.C(c).Where("id in ?", dict.Children.ToSlice()).Order("sort asc,id asc").Find(&dicts); result.Error != nil {
		if xdb.IsNotFound(result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return dicts, nil
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/dict/list", run())
}
