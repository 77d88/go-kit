package pro

import (
	"github.com/77d88/go-kit/plugins/xdb"
)

const TableNameMenu = "s_sys_menu"

func init() {
	xdb.AddModels(&Menu{})
}

// Menu 前端路由和菜单
type Menu struct {
	xdb.BaseModel
	UpdateUser    int64          `gorm:"comment:更新人"`
	Path          string         `gorm:"comment:更新人"`
	ComponentPath string         `gorm:"comment:组件路径"`
	Redirect      string         `gorm:"comment:跳转地址"`
	Name          string         `gorm:"comment:名称"`
	NameZh        string         `gorm:"comment:中文名称"`
	MataTitle     string         `gorm:"comment:页面标题"`
	MataKeywords  string         `gorm:"comment:页面关键字"`
	MetaIcon      string         `gorm:"comment:图标"`
	MetaHide      bool           `gorm:"comment:隐藏"`
	Sort          int            `gorm:"comment:排序"`
	MetaNoLevel   bool           `gorm:"comment:无级"`      // 不自动提升等级 默认情况下只有一个子集菜单时自动升级为上级菜单
	Permission    string         `gorm:"comment:权限"`      // 权限 对应权限表的code
	RouteParams   string         `gorm:"comment:路由访问时参数"` //
	Children      *xdb.Int8Array `gorm:"comment:子菜单"`
	IsSystem      bool           `gorm:"comment:系统菜单"` // 系统菜单只能手动调整
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
			ComponentPath: "common/layout/index",
			Path:          "/superadmin/manager",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
			Children:      xdb.NewInt8Array(2, 3, 4, 5, 6),
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(2),
			Name:          "superadmin-index",
			NameZh:        "系统管理",
			Path:          "/superadmin/index",
			ComponentPath: "system/super/index",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(3),
			Name:          "superadmin-menu",
			NameZh:        "菜单管理",
			Path:          "/superadmin/menu",
			ComponentPath: "system/menu/index",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(4),
			Name:          "superadmin-dict",
			NameZh:        "字典管理",
			Path:          "/superadmin/dict",
			ComponentPath: "system/dict/index",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(5),
			Name:          "superadmin-test-index",
			NameZh:        "测试表格",
			Path:          "/superadmin/test_index",
			ComponentPath: "test/list/index",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
		},
		&Menu{
			BaseModel:     xdb.NewBaseModel(6),
			Name:          "superadmin-test-form",
			NameZh:        "测试表单",
			Path:          "/superadmin/test_form",
			ComponentPath: "test/form/index",
			Permission:    xdb.NewTextArray(RoleSuperAdmin),
			IsSystem:      true,
		},
	}
}
