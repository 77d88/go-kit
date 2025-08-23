package pro

import (
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
)

const TableNameMenu = "s_sys_menu"

//sql:
// INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (1, '2025-05-23 05:24:12.651612+00', '2025-08-21 05:35:17.047207+00', NULL, -1, '/superadmin/manager', 'Layout', '/superadmin/index', 'superadmin-manager', '超级系统管理', '', '', '', 'f', 0, 'f', '{superAdmin}', '', '{2,3,4,5,6,7}', 't', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (4, '2025-05-23 05:24:12.757264+00', '2025-05-23 05:24:12.757264+00', NULL, 0, '/superadmin/dict', 'pro/dict/index', '', 'superadmin-dict', '字典管理', '', '', '', 'f', 0, 'f', '{superAdmin}', '', NULL, 'f', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (3, '2025-05-23 05:24:12.722352+00', '2025-05-23 05:24:12.722352+00', NULL, 0, '/superadmin/menu', 'pro/menu/index', '', 'superadmin-menu', '菜单管理', '', '', '', 'f', 0, 'f', '{superAdmin}', '', NULL, 'f', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (7, '2025-08-21 05:30:10.710795+00', '2025-08-21 06:42:56.15772+00', NULL, -1, '/system/user/index', 'pro/user/index', '', 'superadmin-user', '系统用户', '', '', '', 't', 0, 'f', '{system_user}', '', NULL, 'f', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (2, '2025-05-23 05:24:12.687302+00', '2025-08-21 05:01:14.756478+00', NULL, -1, '/superadmin/index', 'pro/manager/index', '', 'superadmin-index2', '系统管理', '', '', '', 'f', 0, 'f', '{superAdmin}', '', NULL, 'f', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (5, '2025-05-23 05:24:12.794021+00', '2025-08-21 05:00:01.114187+00', NULL, -1, '/superadmin/test_index', 'pro/test/index', '', 'superadmin-test-index', '测试首页', '', '', '2', 'f', 0, 'f', '{superAdmin}', '', NULL, 'f', NULL);
//INSERT INTO "public"."s_sys_menu" ("id", "created_time", "updated_time", "deleted_time", "update_user", "path", "component_path", "redirect", "name", "name_zh", "mata_title", "mata_keywords", "meta_icon", "meta_hide", "sort", "meta_no_level", "permission", "route_params", "children", "root_menu", "layout_component_path") VALUES (6, '2025-05-23 05:24:12.829402+00', '2025-08-21 05:03:36.703834+00', NULL, -1, '/superadmin/test_table', 'pro/test/table/index', '', 'superadmin-test-form', '测试表格', '', '', '', 'f', 0, 'f', '{superAdmin}', '', NULL, 'f', NULL);

// Menu 前端路由和菜单
type Menu struct {
	xpg.BaseModel
	UpdateUser    int64            `gorm:"comment:更新人" json:"updateUser,omitempty"`
	Path          string           `gorm:"comment:更新人" json:"path,omitempty"`
	ComponentPath string           `gorm:"comment:组件路径" json:"componentPath,omitempty"`
	Redirect      string           `gorm:"comment:跳转地址" json:"redirect,omitempty"`
	Name          string           `gorm:"comment:名称" json:"name,omitempty"`
	NameZh        string           `gorm:"comment:中文名称" json:"nameZh,omitempty"`
	MataTitle     string           `gorm:"comment:页面标题" json:"mataTitle,omitempty"`
	MataKeywords  string           `gorm:"comment:页面关键字" json:"mataKeywords,omitempty"`
	MetaIcon      string           `gorm:"comment:图标" json:"metaIcon,omitempty"`
	MetaHide      bool             `gorm:"comment:隐藏" json:"metaHide,omitempty"`
	Sort          int              `gorm:"comment:排序" json:"sort,omitempty"`
	MetaNoLevel   bool             `gorm:"comment:无级" json:"metaNoLevel,omitempty"`      // 不自动提升等级 默认情况下只有一个子集菜单时自动升级为上级菜单
	Permission    []string         `gorm:"comment:权限" json:"permission,omitempty"`       // 权限 对应权限表的code
	RouteParams   string           `gorm:"comment:路由访问时参数" json:"routeParams,omitempty"` //
	Children      xtype.Int64Array `gorm:"comment:子菜单" json:"children,omitempty"`
	RootMenu      bool             `gorm:"comment:系统根目录" json:"rootMenu,omitempty"` // 系统菜单只能手动调整
}

// TableName Res's table name
func (*Menu) TableName() string {
	return TableNameMenu
}
