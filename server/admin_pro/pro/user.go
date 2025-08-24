package pro

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
)

const TableNameUser = "s_sys_user"

func init() {
	//xpg.AddModels(&User{})
}

// User 用户表
type User struct {
	xpg.BaseModel
	UpdateUser          int64            `gorm:"comment:更新人"`
	Password            string           `gorm:"comment:后台登录密码" json:"password"` // 后台登录密码
	Disabled            bool             `gorm:"comment:是否禁用" json:"disabled"`   // 是否禁用
	Username            string           `gorm:"comment:后台登录名称" json:"username"` // 后台登录名称
	Nickname            string           `gorm:"comment:后台显示名称" json:"nickname"` // 后台显示名称
	Avatar              xtype.Int64Array `gorm:"comment:头像" json:"avatar"`       // 头像
	Roles               xtype.Int64Array `gorm:"comment:系统角色" json:"roles"`      // 系统角色
	Email               string           `json:"email"`
	IsReLogin           bool             `json:"isReLogin"`
	ReLoginDesc         string           `json:"reLoginDesc"`
	RoleNames           []string         `json:"roleNames"`
	RolePermissionCodes []string         `json:"RolePermissionCodes"` // 冗余 集合Roles里面的所有Permission Code
}

// TableName Res's table name
func (*User) TableName() string {
	return TableNameUser
}

// AllPermissionCode 获取所有权限码 去重
func (d *User) AllPermissionCode() []string {
	return d.RolePermissionCodes
}

func (d *User) HasPermission(code ...string) bool {
	return xarray.ContainAny(d.AllPermissionCode(), code)
}
