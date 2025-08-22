package items

import (
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xcache"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 前端获取选项
type request struct {
	Name string `json:"name"`
}
type response struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Name string `json:"name"`
}

func handler(c *xhs.Ctx, r *request) (interface{}, error) {
	return xcache.Once[any]("pro.dict.items", time.Minute*10, func() (any, error) {
		var dict pro.Dict
		if result := xdb.C(c).Where("name = ?", r.Name).Find(&dict); result.Error != nil {
			return nil, result.Error
		}
		if dict.Children.IsEmpty() {
			return []struct{}{}, nil
		}
		var items []pro.Dict
		if result := xdb.C(c).Where("id in ?", dict.Children.ToSlice()).Find(&items); result.Error != nil {
			return nil, result.Error
		}
		return xarray.Map(items, func(index int, item pro.Dict) response {
			return response{
				Code: item.Code,
				Desc: item.Desc,
				Name: item.Name,
			}
		}), nil
	})

}

func run() xhs.Handler {
	return xhs.DefaultShouldHandler(handler)
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/dict/items", run(), auth.ForceAuth)
}
