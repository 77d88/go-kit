package pro

import (
	"github.com/77d88/go-kit/plugins/xdb"
)

const TableNameDict = "s_sys_dict"

func init() {
	xdb.AddModels(&Dict{})
}

// Dict 字典表
type Dict struct {
	xdb.BaseModel        // 这个表都物理删除哦
	UpdateUser    int64  `gorm:"comment:更新人"`
	TypeId        int64  `gorm:"comment:字典分类;index"` //对应字典类型ID
	Val           int    `gorm:"comment:字典值"`
	Desc          string `gorm:"comment:字典描述"`
	Sort          int    `gorm:"comment:字典排序"`
	IsType        bool   `gorm:"comment:是否字典类型;index"`
}

// TableName Res's table name
func (*Dict) TableName() string {
	return TableNameDict
}

func (d *Dict) ToResponse() *DictDst {
	return &DictDst{
		Id:   d.ID,
		Desc: d.Desc,
		Sort: d.Sort,
		Val:  d.Val,
		Type: d.TypeId,
	}
}

type DictDst struct {
	Id   int64  `json:"id,string"`
	Val  int    `json:"val"`
	Desc string `json:"desc"`
	Sort int    `json:"sort"`
	Type int64  `json:"type,string"`
}

func ToDictResponses(dicts []*Dict) []*DictDst {
	res := make([]*DictDst, 0)
	if len(dicts) == 0 {
		return res
	}

	for _, dict := range dicts {
		res = append(res, dict.ToResponse())
	}
	return res
}
