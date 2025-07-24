package api_menu

import (
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
	"testing"
)

func TestName(t *testing.T) {

	menus := []*pro.Menu{
		{
			BaseModel: xdb.NewBaseModel(0),
			Path:      "/index",
			Name:      "index",
			NameZh:    "欢迎",
			Sort:      3,
		},
		{
			BaseModel: xdb.NewBaseModel(1),
			Path:      "/system",
			Name:      "System",
			NameZh:    "系统管理",
			Children:  xdb.NewInt8Array(2, 3),
			Sort:      2,
		},
		{
			BaseModel: xdb.NewBaseModel(2),
			Path:      "/system/user",
			Name:      "User",
			NameZh:    "用户管理",
			Children:  xdb.NewInt8Array(4),
		},
		{
			BaseModel: xdb.NewBaseModel(3),
			Path:      "/system/role",
			Name:      "Role",
			NameZh:    "角色管理",
		},
		{
			BaseModel:  xdb.NewBaseModel(4),
			Path:       "/system/user/list",
			Name:       "UserList",
			NameZh:     "用户列表",
			Permission: xdb.NewTextArray("index2"),
		},
	}

	routers := ConvertMenusToRouter(menus, "index")
	jsonData, _ := json.MarshalIndent(routers, "", "  ")
	fmt.Println(string(jsonData))

}
