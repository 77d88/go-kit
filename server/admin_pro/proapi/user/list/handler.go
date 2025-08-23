package list

import (
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户列表
type response struct {
	Id          int64            `json:"id,string"`
	Password    string           `json:"password,omitempty"`
	Disabled    bool             `json:"disabled"`
	Username    string           `json:"username"`
	Nickname    string           `json:"nickname"`
	Avatar      xtype.Int64Array `json:"avatar"`
	Roles       xtype.Int64Array `json:"roles"`
	Permission  xtype.Int64Array `json:"permission"`
	Email       string           `json:"email"`
	IsReLogin   bool             `json:"isReLogin"`
	ReLoginDesc string           `json:"reLoginDesc"`
	UpdatedTime time.Time        `json:"updatedTime"`
}

type request struct {
	Page     xdb.PageSearch `json:"page"`
	Name     string         `json:"name"`
	Disabled *bool          `json:"disabled,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	var users []pro.User
	result := xpg.C(c).Model(&pro.User{}).Order("id desc").
		XWhere(r.Name != "", "(username ilike @name or nickname ilike @name)", xdb.Param("name", xdb.WarpLike(r.Name))).
		XWhere(r.Disabled != nil, "disabled = ?", r.Disabled).FindPage(&users, r.Page, true)

	if result.Error != nil {
		return nil, result.Error
	}
	return xhs.NewResp(xarray.Map(users, func(i int, d pro.User) *response {
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
			UpdatedTime: d.UpdatedTime,
		}
	}), result.Total), nil

}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/list", run(), auth.ForceAuth)
}
