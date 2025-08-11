package pro

import (
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

const TableNamePermission = "s_sys_permission"

func init() {
	xdb.AddModels(&Permission{})
}

// Permission 权限配置
type Permission struct {
	xdb.BaseModel
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

func (d *Permission) InitData() []xdb.GromModel {
	return []xdb.GromModel{
		&Permission{
			BaseModel: xdb.NewBaseModel(1),
			Code:      "sys.user.list",
			Desc:      "用户列表",
		},
		&Permission{
			BaseModel: xdb.NewBaseModel(2),
			Code:      "sys.role.list",
			Desc:      "角色列表",
		},
		&Permission{
			BaseModel: xdb.NewBaseModel(3),
			Code:      "pro.home",
			Desc:      "首页",
		},
	}
}
