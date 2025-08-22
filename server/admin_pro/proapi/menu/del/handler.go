package del

import (
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

// 菜单删除
type response struct {
}

type request struct {
	Id int64 `json:"id,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Id <= 0 {
		return nil, xerror.New("参数错误")
	}

	if result := xdb.C(c).Where("id = ?", r.Id).Updates(map[string]interface{}{
		"deleted_time": time.Now(),
		"update_user":  c.GetUserId(),
	}); result.Error != nil {
		return nil, result.Error
	}

	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/menu/del", run(), auth.ForceAuth)
}
