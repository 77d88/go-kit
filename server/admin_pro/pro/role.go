package pro

import (
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
)

const TableNameRole = "s_sys_role"

// Role 权限配置
type Role struct {
	xpg.BaseModel
	UpdateUser      int64            `gorm:"comment:更新人"`
	Name            string           `gorm:"comment:权限码"` //全局唯一
	Permission      xtype.Int64Array `gorm:"comment:权限集合"`
	PermissionCodes []string         `gorm:"comment:权限码"` // 冗余字段
}

// TableName Res's table name
func (*Role) TableName() string {
	return TableNameRole
}

type RoleDst struct {
	Id         any              `json:"id,string" `
	Name       string           `json:"name"`
	Permission xtype.Int64Array `json:"permission"`
}

func (d *Role) ToResponse() *RoleDst {
	return &RoleDst{
		Id:         d.ID,
		Name:       d.Name,
		Permission: d.Permission,
	}
}
