package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
	"gorm.io/gorm"
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
	ParentId int64  `json:"parentId,string"`
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
		if result := xdb.C(c).Where("name = ?", r.Name).Find(&dict); result.Error != nil {
			return nil, result.Error
		}
		if dict.ID > 0 && dict.ID != r.Id {
			return nil, xerror.New("字典名称不能重复")
		}
	}

	var parent pro.Dict
	if r.ParentId > 0 {
		if result := xdb.C(c).Where("id = ?", r.ParentId).Find(&parent); result.Error != nil {
			return nil, result.Error
		}
	}

	err = xdb.C(c).Transaction(func(tx *gorm.DB) error {
		result := xdb.SaveMap[pro.Dict](tx, r, map[string]interface{}{
			"update_user": c.GetUserId(),
			"ParentId":    xdb.ToMapIgnore,
		})
		if result.Error != nil {
			return result.Error
		}
		if r.ParentId > 0 {
			// 如果不包含这个子项则新增
			children := parent.Children
			if children == nil {
				children = &xdb.Int8Array{}
			}
			if !children.Contain(result.RowId) {
				children.AppendIfNotExist(result.RowId)
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
