package pro

import (
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
)

const TableNamePermission = "s_sys_permission"

// Permission 权限配置
type Permission struct {
	xpg.BaseModel
	UpdateUser int64  `gorm:"comment:更新人"`
	Code       string `gorm:"comment:权限码;index"` //全局唯一
	Desc       string `gorm:"comment:权限描述"`
}

// TableName Res's table name
func (*Permission) TableName() string {
	return TableNamePermission
}

type PermissionDst struct {
	Id   int64  `json:"id,string"`
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func (d *Permission) ToResponse() *PermissionDst {
	return &PermissionDst{
		Id:   d.ID,
		Code: d.Code,
		Desc: d.Desc,
	}
}
