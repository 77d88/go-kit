package pro

import (
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

const TableNameMenu = "s_sys_menu"

func init() {
	xdb.AddModels(&Menu{})
}

// Menu 前端路由和菜单
type Menu struct {
	xdb.BaseModel
	UpdateUser    int64          `gorm:"comment:更新人" json:"updateUser,omitempty"`
	Path          string         `gorm:"comment:更新人" json:"path,omitempty"`
	ComponentPath string         `gorm:"comment:组件路径" json:"componentPath,omitempty"`
	Redirect      string         `gorm:"comment:跳转地址" json:"redirect,omitempty"`
	Name          string         `gorm:"comment:名称" json:"name,omitempty"`
	NameZh        string         `gorm:"comment:中文名称" json:"nameZh,omitempty"`
	MataTitle     string         `gorm:"comment:页面标题" json:"mataTitle,omitempty"`
	MataKeywords  string         `gorm:"comment:页面关键字" json:"mataKeywords,omitempty"`
	MetaIcon      string         `gorm:"comment:图标" json:"metaIcon,omitempty"`
	MetaHide      bool           `gorm:"comment:隐藏" json:"metaHide,omitempty"`
	Sort          int            `gorm:"comment:排序" json:"sort,omitempty"`
	MetaNoLevel   bool           `gorm:"comment:无级" json:"metaNoLevel,omitempty"`      // 不自动提升等级 默认情况下只有一个子集菜单时自动升级为上级菜单
	Permission    *xdb.TextArray `gorm:"comment:权限" json:"permission,omitempty"`       // 权限 对应权限表的code
	RouteParams   string         `gorm:"comment:路由访问时参数" json:"routeParams,omitempty"` //
	Children      *xdb.Int8Array `gorm:"comment:子菜单" json:"children,omitempty"`
	RootMenu      bool           `gorm:"comment:系统根目录" json:"rootMenu,omitempty"` // 系统菜单只能手动调整
}

// TableName Res's table name
func (*Menu) TableName() string {
	return TableNameMenu
}

func (m *Menu) InitData() []xdb.GromModel {
	return []xdb.GromModel{
		&Menu{
			BaseModel:     xdb.NewBaseModel(1),
			Name:          "superadmin-manager",
			NameZh:        "超级系统管理",
			Redirect:      "/superadmin/index",
			ComponentPath: "Layout",
			Path:          "/superadmin/manager",
			Permission:    xdb.NewTextArray("superadmin_manager"),
			RootMenu:      true,
			Children:      xdb.NewInt8Array(2, 3, 4, 5, 6, 7, 8, 9),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(2),
			Name:          "superadmin-index",
			NameZh:        "系统管理",
			Path:          "/superadmin/index",
			ComponentPath: "system/super/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(3),
			Name:          "superadmin-menu",
			NameZh:        "菜单管理",
			Path:          "/superadmin/menu",
			ComponentPath: "system/menu/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(4),
			Name:          "superadmin-menu",
			NameZh:        "系统用户",
			Path:          "/superadmin/user",
			ComponentPath: "system/user/index",
			Permission:    xdb.NewTextArray("superadmin_user"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(5),
			Name:          "superadmin-dict",
			NameZh:        "字典管理",
			Path:          "/superadmin/dict",
			ComponentPath: "system/dict/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(6),
			Name:          "superadmin-role",
			NameZh:        "角色管理",
			Path:          "/superadmin/role",
			ComponentPath: "system/role/index",
			Permission:    xdb.NewTextArray("superadmin_role"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(7),
			Name:          "superadmin-permission",
			NameZh:        "权限管理",
			Path:          "/superadmin/permission",
			ComponentPath: "system/permission/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(8),
			Name:          "superadmin-test-index",
			NameZh:        "测试表格",
			Path:          "/superadmin/test_index",
			ComponentPath: "test/list/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(9),
			Name:          "superadmin-test-form",
			NameZh:        "测试表单",
			Path:          "/superadmin/test_form",
			ComponentPath: "test/form/index",
			Permission:    xdb.NewTextArray("superAdmin"),
		},
	}
}
