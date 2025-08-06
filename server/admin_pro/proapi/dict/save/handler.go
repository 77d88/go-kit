package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典保存
type response struct {
}

type request struct {
	Id     int64  `json:"id,string"`
	Val    int    `json:"val"`
	Desc   string `json:"desc"`
	Sort   int    `json:"sort"`
	TypeId int64  `json:"typeId,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Id > 0 {
		dict := pro.Dict{}
		if result := xdb.Ctx(c).WithId(r.Id).Find(&dict); result.Error != nil {
			return nil, result.Error
		}
		if !dict.IsType {
			take := xdb.Ctx(c).Where("val = ? and id != ?", dict.Val, r.Id).Take(&pro.Dict{})
			if !take.IsNotFound() {
				return nil, xerror.New("字典类型已存在")
			}
		}
		if result := xdb.Ctx(c).Model(&pro.Dict{}).WithId(r.Id).Updates(map[string]any{
			"desc":        r.Desc,
			"sort":        r.Sort,
			"val":         r.Val,
			"type_id":     r.TypeId,
			"update_user": c.GetUserId(),
		}); result.Error != nil {
			return nil, result.Error
		}
	} else {
		if r.Val <= 0 {
			return nil, xerror.New("参数错误")
		}
		if result := xdb.Ctx(c).Model(&pro.Dict{}).Create(&pro.Dict{
			Desc:       r.Desc,
			Sort:       r.Sort,
			TypeId:     r.TypeId,
			UpdateUser: c.GetUserId(),
			Val:        r.Val,
		}); result.Error != nil {
			return nil, result.Error
		}
	}

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
