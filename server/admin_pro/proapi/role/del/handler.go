package del

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
	"time"
)

// 角色删除
type response struct {
}

type request struct {
	Id int64 `json:"id,omitempty"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Id <= 0 {
		return nil, xerror.New("参数错误")
	}
	var role pro.Role
	if result := xdb.Ctx(c).WithId(r.Id).Take(&role); result.Error != nil {
		return nil, result.Error
	}
	var count int64
	if result := xdb.Ctx(c).Model(&pro.Role{}).Where("deleted_time is null and roles && ?", xdb.NewInt8Array(r.Id)).Count(&count); result.Error != nil {
		return nil, result.Error
	}
	if count > 0 {
		return nil, xerror.New("角色已分配给用户，不能删除")
	}
	if result := xdb.Ctx(c).Model(&pro.Role{}).WithId(r.Id).
		Updates(map[string]interface{}{
			"update_user":  c.GetUserId(),
			"deleted_time": time.Now(),
		}); result.Error != nil {
		return nil, result.Error
	}

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
