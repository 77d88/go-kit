package list

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典列表
type response struct {
}

type request struct {
	TypeId int64 `json:"typeId,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.TypeId <= 0 {
		return nil, xerror.New("参数错误")
	}
	var dict []*pro.Dict
	if result := xdb.Ctx(c).Where("not is_type and type = ?", r.TypeId).Order("sort asc").Find(&dict); result.Error != nil {
		return nil, result.Error
	}
	return dict, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
