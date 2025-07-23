package admin_pro

import (
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_db"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_dict"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_menu"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_permission"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_role"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro/api_user"
	"github.com/77d88/go-kit/plugins/xe"
)

func RegisterApi(api *xe.Engine) {
	api_db.Register(api, "/v1/pro/sys/db")
	api_menu.Register(api, "/v1/pro/sys/menu")
	api_dict.Register(api, "/v1/pro/sys/dict")
	api_user.Register(api, "/v1/pro/sys/user")
	api_permission.Register(api, "/v1/pro/sys/permission")
	api_role.Register(api, "/v1/pro/sys/role")
}
