package list

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户列表
type response struct {
	Id          int64          `json:"id,string"`
	Password    string         `json:"password,omitempty"`
	Disabled    bool           `json:"disabled"`
	Username    string         `json:"username"`
	Nickname    string         `json:"nickname"`
	Avatar      *xdb.Int8Array `json:"avatar"`
	Roles       *xdb.Int8Array `json:"roles"`
	Permission  *xdb.Int8Array `json:"permission"`
	Email       string         `json:"email"`
	IsReLogin   bool           `json:"isReLogin"`
	ReLoginDesc string         `json:"reLoginDesc"`
}

type request struct {
	xdb.PageSearch
	Name     string `json:"name"`
	Disabled *bool  `json:"disabled"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {

	var (
		total int64
		users []*pro.User
	)
	if result := xdb.Ctx(c).Model(&pro.User{}).
		Where("not is_super_admin"). // 超级管理员不显示出来
		XWhere(r.Name != "", "username ilike @name || nickname ilike @name", xdb.Param("name", xdb.WarpLike(r.Name))).
		XWhere(r.Disabled != nil, "disabled = ?", r.Disabled).
		IdDesc().FindPage(r, &users, &total); result.Error != nil {
		return nil, result.Error
	}
	return xhs.NewResp(xarray.Map(users, func(i int, d *pro.User) *response {
		return &response{
			Id:          d.ID,
			Disabled:    d.Disabled,
			Nickname:    d.Nickname,
			Username:    d.Username,
			Avatar:      d.Avatar,
			Roles:       d.Roles,
			Permission:  d.Permission,
			IsReLogin:   d.IsReLogin,
			ReLoginDesc: d.ReLoginDesc,
		}
	})), nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
