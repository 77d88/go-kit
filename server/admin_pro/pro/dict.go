package pro

import (
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

const TableNameDict = "s_sys_dict"

func init() {
	xdb.AddModels(&Dict{})
}

// Dict 字典表
type Dict struct {
	xdb.BaseModel                // 这个表都物理删除哦
	UpdateUser    int64          `gorm:"comment:更新人" json:"updateUser,omitempty"`
	Code          int            `gorm:"comment:字典值" json:"code,omitempty"`
	Desc          string         `gorm:"comment:字典描述" json:"desc,omitempty"`
	Name          string         `gorm:"comment:字典名称" json:"name,omitempty"`
	Sort          int            `gorm:"comment:字典排序" json:"sort,omitempty"`
	Children      *xdb.Int8Array `gorm:"comment:子菜单" json:"children,omitempty" `
	Root          bool           `gorm:"comment:字典类型" json:"root,omitempty" `
}

// TableName Res's table name
func (*Dict) TableName() string {
	return TableNameDict
}
