package xdb

import (
	"encoding/json"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/jackc/pgtype"
)

func init() {
}

type TextArray struct {
	pgtype.TextArray
}

func (i *TextArray) IsEmpty() bool {
	return i == nil || len(i.Elements) == 0
}

func (i *TextArray) ToSlice() []string {
	var goArray = make([]string, 0)
	if i == nil {
		return goArray
	}
	for _, elem := range i.Elements {
		if elem.Status == pgtype.Present {
			val := elem.String
			goArray = append(goArray, val)
		}
	}
	return goArray
}

func NewTextArray(s ...string) *TextArray {
	// 初始化一个空的 pgtype.Int8Array
	var pgArray = TextArray{}
	// 尝试将 int64 值切片设置到 pgtype.Int8Array
	err := pgArray.Set(s)
	// 如果在设置过程中发生错误，函数仍会返回初始化的 pgtype.Int8Array，但其值可能不符合预期。
	if err != nil {
		return &pgArray
	}
	// 如果设置过程成功，返回初始化的 pgtype.Int8Array
	return &pgArray
}

func NewTextArrayUnique(s ...string) *TextArray {
	return NewTextArray(xarray.Unique(s)...)
}

func NewTextArrayByPointer(s ...*string) *Int8Array {
	var pgArray = Int8Array{}
	err := pgArray.Set(s)
	if err != nil {
		return &pgArray
	}
	return &pgArray
}

// Append 追加元素
func (i *TextArray) Append(is ...string) {
	if len(is) == 0 {
		return
	}

	slice := i.ToSlice()
	slice = append(slice, is...)
	i.Set(slice)
}

// AppendIfNotExist 追加元素，如果元素不存在
// 返回是否追加成功
func (i *TextArray) AppendIfNotExist(is ...string) bool {
	flag := false
	if len(is) == 0 {
		return false
	}
	slice := i.ToSlice()
	for _, x := range is {
		if xarray.Contain(slice, x) {
			continue
		}
		slice = append(slice, x)
		flag = true
	}
	if flag {
		err := i.Set(slice)
		if err != nil {
			xlog.Errorf(nil, "AppendIfNotExist error: %v", err)
			return false
		}
	}
	return flag
}

func (i *TextArray) Contain(val string) bool {
	return xarray.Contain(i.ToSlice(), val)
}
func (i *TextArray) ContainBy(predicate func(item string) bool) bool {
	return xarray.ContainBy(i.ToSlice(), predicate)
}

// MarshalJSON 实现 json.Marshaler 接口
func (i *TextArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ToSlice())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *TextArray) UnmarshalJSON(data []byte) error {
	var strArray []string
	if err := json.Unmarshal(data, &strArray); err != nil {
		return err
	}
	return i.Set(strArray)
}

func (*TextArray) GormDataType() string {
	return "text[]"
}
