package types

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典类型查询
type response struct {
	Id   int64  `json:"id,string"`
	Desc string `json:"desc"`
	Val  int    `json:"val"`
	Sort int    `json:"sort"`
}

type request struct {
	xdb.PageSearch
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if result := xdb.FindPage[pro.Dict](xdb.C(c).Where("is_type").Order("sort asc"), r, true); result.Error != nil {
		return nil, result.Error
	} else {
		return xhs.NewResp(xarray.Map(result.List, func(i int, item pro.Dict) *response {
			return &response{
				Id:   item.ID,
				Desc: item.Desc,
				Val:  item.Val,
				Sort: item.Sort,
			}
		}), result.Total), nil
	}

}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
