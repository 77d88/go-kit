package api_dict

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	pro2 "github.com/77d88/go-kit/server/admin_pro/pro"
	"time"
)

func getTypes(c *xhs.Ctx) {
	var dict []*pro2.Dict
	c.Fatalf(xdb.Ctx(c).Where("is_type").Order("sort asc").Find(&dict))
	c.Send(pro2.ToDictResponses(dict))
}

type typeListReq struct {
	TypeId int `json:"typeId"`
}

func typeList(c *xhs.Ctx) {
	var req typeListReq
	c.ShouldBind(&req)
	c.Fatalf(req.TypeId <= 0, "参数错误")
	var dict []*pro2.Dict
	c.Fatalf(xdb.Ctx(c).Where("not is_type and type = ?", req.TypeId).Order("sort asc").Find(&dict))
	c.Send(pro2.ToDictResponses(dict))
}
func deleteDict(c *xhs.Ctx) {
	var req pro2.Dict
	c.ShouldBind(&req)
	c.Fatalf(req.ID <= 0, "参数错误")

	var dict pro2.Dict
	c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Where("id = ?", req.ID).First(&dict))

	c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Where("id = ?", dict.ID).Updates(map[string]any{
		"deleted_time": time.Now(),
		"update_user":  c.GetUserIdAssert(),
	}))

	if dict.IsType { // 删除字典类型相关的所有字典
		c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Where("type = ?", dict.Type).Updates(map[string]any{
			"deleted_time": time.Now(),
			"update_user":  c.GetUserIdAssert(),
		}))
	}

}
func saveDict(c *xhs.Ctx) {
	var req pro2.DictDst
	c.ShouldBind(&req)
	if req.Id > 0 {
		var dict pro2.Dict
		c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Where("id = ?", req.Id).First(&dict))

		if !dict.IsType {
			r := xdb.Ctx(c).Model(&pro2.Dict{}).Where("type = ? and val = ?", dict.Type, req.Val).First(&pro2.Dict{})
			if r.IsNotFound() {
				c.SendError(xerror.New("字典已存在"))
				return
			}
			c.Fatalf(r)
		}

		c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Where("id = ?", req.Id).Updates(map[string]any{
			"desc":        req.Desc,
			"sort":        req.Sort,
			"val":         req.Val,
			"update_user": c.GetUserIdAssert(),
		}))
	} else {
		c.Fatalf(req.Type <= 0, "参数错误")
		c.Fatalf(req.Val <= 0, "参数错误")
		r := xdb.Ctx(c).Model(&pro2.Dict{}).Where("type = ? and val = ?", req.Type, req.Val).First(&pro2.Dict{})
		if r.IsNotFound() {
			c.SendError(xerror.New("字典已存在"))
			return
		}
		c.Fatalf(r)

		c.Fatalf(xdb.Ctx(c).Model(&pro2.Dict{}).Create(&pro2.Dict{
			Desc:       req.Desc,
			Sort:       req.Sort,
			Type:       req.Type,
			UpdateUser: c.GetUserIdAssert(),
			Val:        req.Val,
		}))
	}

}

func Register(api *xe.Engine, path string) {
	api.RegisterPost(path+"/getTypes", pro2.SuperAdmin, getTypes) // 获取字典类型列表
	api.RegisterPost(path+"/list", pro2.SuperAdmin, typeList)     // 获取类型的字典列表
	api.RegisterPost(path+"/delete", pro2.SuperAdmin, deleteDict) // 删除字典
	api.RegisterPost(path+"/save", pro2.SuperAdmin, saveDict)
}

func Default(api *xe.Engine) {
	Register(api, "/x/sys/dict")
}
