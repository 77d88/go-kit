package del

import (
	"time"

	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 字典删除
type response struct {
}

type request struct {
	ID int64 `json:"id,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	var dict pro.Dict
	if result := xpg.C(c).Model(&pro.Dict{}).Where("id = ?", r.ID).First(&dict); result.Error != nil {
		return nil, err
	}

	err = xpg.C(c).Transaction(func(tx *xpg.Inst) error {
		if result := tx.Model(&pro.Dict{}).Where("id = ?", dict.ID).Updates(map[string]any{
			"deleted_time": time.Now(),
			"update_user":  c.GetUserId(),
		}); result.Error != nil {
			return result.Error
		}
		xlog.Warnf(c, "删除字典成功: %d, %s, %d, %s", dict.ID, dict.Name, dict.Code, dict.Desc)
		if !dict.Children.IsEmpty() {
			slice := dict.Children.ToSlice()
			if result := tx.Model(&pro.Dict{}).Where("id = ANY(?)", slice).Updates(map[string]any{
				"deleted_time": time.Now(),
				"update_user":  c.GetUserId(),
			}); result.Error != nil {
				return result.Error
			}
			xlog.Warnf(c, "相关字典删除完成: %v", slice)
		}
		return nil

	})
	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/dict/del", run(), auth.ForceAuth)
}
