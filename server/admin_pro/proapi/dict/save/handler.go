package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典保存
type response struct {
}

type request struct {
	Id       int64  `json:"id,string"`
	Code     int    `json:"code,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Name     string `json:"name,omitempty"`
	Sort     int    `json:"sort,omitempty"`
	Root     bool   `json:"root,omitempty" `
	ParentId int64  `json:"parentId,string" db:"-"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Name == "" {
		return nil, xerror.New("参数错误:名称不能为空")
	}
	if r.Code < 0 {
		return nil, xerror.New("参数错误:字典值不能小于0")
	}

	// 根节点name不能重复
	if r.ParentId == 0 {
		var dict pro.Dict
		if result := xpg.C(c).Where("name = ?", r.Name).Find(&dict); result.Error != nil {
			return nil, result.Error
		}
		if dict.ID > 0 && dict.ID != r.Id {
			return nil, xerror.New("字典名称不能重复")
		}
	}

	var parent pro.Dict
	if r.ParentId > 0 {
		if result := xpg.C(c).Where("id = ?", r.ParentId).Find(&parent); result.Error != nil {
			return nil, result.Error
		}
	}

	err = xpg.C(c).Transaction(func(tx *xpg.Inst) error {
		result := tx.Model(&pro.Dict{}).Save(r, func(m map[string]interface{}) {
			m["update_user"] = c.GetUserId()
		})
		if result.Error != nil {
			return result.Error
		}
		if r.ParentId > 0 {
			// 如果不包含这个子项则新增
			children := parent.Children
			if children == nil {
				children = make(xtype.Int64Array, 0)
			}
			if !children.Contain(result.RowId) {
				children = append(children, result.RowId)
				if result := tx.Model(&parent).Where("id = ?", parent.ID).Update("children", children); result.Error != nil {
					return result.Error
				}
			}

		}
		return nil
	})

	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/dict/save", run(), auth.ForceAuth)
}
