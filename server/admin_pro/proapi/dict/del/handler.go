package del

import (
	"time"

	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
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
	if result := xdb.C(c).Model(&pro.Dict{}).Where("id = ?", r.ID).First(&dict); result.Error != nil {
		return nil, err
	}

	if result := xdb.C(c).Model(&pro.Dict{}).Where("id = ?", dict.ID).Updates(map[string]any{
		"deleted_time": time.Now(),
		"update_user":  c.GetUserId(),
	}); result.Error != nil {
		return nil, result.Error
	}
	xlog.Warnf(c, "删除字典成功: %d, %d, %d, %s", dict.ID, dict.TypeId, dict.Val, dict.Desc)

	if dict.IsType { // 删除字典类型相关的所有字典
		if result := xdb.C(c).Model(&pro.Dict{}).Where("type_id = ?", dict.ID).Updates(map[string]any{
			"deleted_time": time.Now(),
			"update_user":  c.GetUserId(),
		}); result.Error != nil {
			return nil, result.Error
		}
		xlog.Warnf(c, "删除字典类型成功: %d,%d, %s", dict.ID, dict.Val, dict.Desc)
	}

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
