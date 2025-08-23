package info

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 用户信息
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
}

type request struct {
	Id int64 `json:"id,string"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Id <= 0 {
		return nil, xerror.New("参数错误:Id不能为空")
	}
	var user pro.User
	if result := xpg.C(c).Where("id = ?", r.Id).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return &response{
		Id:          user.ID,
		Disabled:    user.Disabled,
		Nickname:    user.Nickname,
		Username:    user.Username,
		Avatar:      user.Avatar,
		Roles:       user.Roles,
		Permission:  user.Permission,
		IsReLogin:   user.IsReLogin,
		ReLoginDesc: user.ReLoginDesc,
	}, nil
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/user/info", run(), auth.ForceAuth)
}
