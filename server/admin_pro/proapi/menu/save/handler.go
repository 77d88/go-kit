package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 菜单保存
type response struct {
}

type request struct {
	Id            int64          `json:"id,string"`
	Path          string         `json:"path"`
	ComponentPath string         `json:"componentPath"`
	Redirect      string         `json:"redirect"`
	Name          string         `json:"name"`
	NameZh        string         `json:"nameZh"`
	MataTitle     string         `json:"mataTitle"`
	MataKeywords  string         `json:"mataKeywords"`
	MetaIcon      string         `json:"metaIcon"`
	MetaHide      bool           `json:"metaHide"`
	Sort          int            `json:"sort"`
	MetaNoLevel   bool           `json:"metaNoLevel"`
	RouteParams   string         `json:"routeParams"`
	RootMenu      bool           `json:"rootMenu"`
	ParentId      int64          `json:"parentId,string"`
	Permission    *xdb.TextArray `json:"permission"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Name == "" {
		return nil, xerror.New("参数错误:名称不能为空")
	}
	var parent pro.Menu
	if r.ParentId > 0 {
		if result := xpg.C(c).Where("id = ?", r.ParentId).Find(&parent); result.Error != nil {
			return nil, result.Error
		}
	}

	err = xpg.C(c).Transaction(func(tx *xpg.Inst) error {
		result := tx.Table(pro.TableNameMenu).Save(r, func(m map[string]interface{}) {
			m["update_user"] = c.GetUserId()
			delete(m, "parent_id")
		})

		if result.Error != nil {
			return result.Error
		}

		if r.ParentId > 0 {
			// 如果不包含这个菜单
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
	xsh.POST("/pro/menu/save", run(), auth.ForceAuth)
}
