package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户保存
type response struct {
}

type request struct {
	Id       int64          `json:"id,string"`
	Password string         `json:"password,omitempty"`
	Disabled bool           `json:"disabled"`
	Username string         `json:"username"`
	Nickname string         `json:"nickname"`
	Avatar   *xdb.Int8Array `json:"avatar"`
	Email    string         `json:"email"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	if r.Username == "" {
		return nil, xerror.New("用户名不能为空")
	}

	if r.Id <= 0 && r.Password == "" {
		return nil, xerror.New("新用户必须要输入密码")
	}

	var user pro.User
	if result := xdb.Ctx(c).Where("username = ?", r.Username).Find(&user); result.Error != nil {
		return nil, result.Error
	}

	if user.ID > 0 && user.ID != r.Id {
		return nil, xerror.New("用户名已存在")
	}
	if result := xdb.Ctx(c).SaveMap(&pro.User{}, r, map[string]interface{}{
		"update_user": c.GetUserId(),
	}); result.Error != nil {
		return nil, result.Error
	}
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
