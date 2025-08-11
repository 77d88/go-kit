package xdb

import (
	"encoding/json"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/jackc/pgtype"
)

// Integer 接口定义了整数的基本操作
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func init() {
}

type Int8Array struct {
	pgtype.Int8Array
}

func (i *Int8Array) IsEmpty() bool {
	return i == nil || len(i.Elements) == 0
}

func (i *Int8Array) ToStrings() []string {
	// 初始化一个空的字符串数组
	var goArray = make([]string, 0)
	if i == nil {
		return goArray
	}

	// 遍历输入数组的每个元素
	for _, elem := range i.Elements {
		if elem.Status == pgtype.Present {
			val := elem.Int
			goArray = append(goArray, xparse.ToString(val, 0))
		}
	}
	// 返回处理后的字符串数组
	return goArray
}

func (i *Int8Array) ToSlice() []int64 {
	var goArray = make([]int64, 0)
	if i == nil {
		return goArray
	}
	for _, elem := range i.Elements {
		if elem.Status == pgtype.Present {
			val := elem.Int
			goArray = append(goArray, val)
		}
	}
	return goArray
}

func NewInt8Array(s ...int64) *Int8Array {
	// 初始化一个空的 pgtype.Int8Array
	var pgArray = Int8Array{}
	// 尝试将 int64 值切片设置到 pgtype.Int8Array
	err := pgArray.Set(s)
	// 如果在设置过程中发生错误，函数仍会返回初始化的 pgtype.Int8Array，但其值可能不符合预期。
	if err != nil {
		return &pgArray
	}
	// 如果设置过程成功，返回初始化的 pgtype.Int8Array
	return &pgArray
}

func NewInt8ArrayUnique(s ...int64) *Int8Array {
	return NewInt8Array(xarray.Unique(s)...)
}
func MergeInt8Array(s ...*Int8Array) *Int8Array {
	var pgArray = Int8Array{}
	for _, item := range s {
		if item == nil {
			continue
		}
		pgArray.Append(item.ToSlice()...)
	}
	return &pgArray
}

func MergeInt8ArrayUnique(s ...*Int8Array) *Int8Array {
	var pgArray = Int8Array{}
	for _, item := range s {
		if item == nil {
			continue
		}
		pgArray.AppendIfNotExist(item.ToSlice()...)
	}
	return &pgArray
}

func NewInt8ArrayByPointer(s ...*int64) *Int8Array {
	var pgArray = Int8Array{}
	err := pgArray.Set(s)
	if err != nil {
		return &pgArray
	}
	return &pgArray
}

func NewInt8ArrayByStrings(s ...string) *Int8Array {
	return NewInt8Array(xarray.MapBy(s, xparse.WarpToMap(xparse.ToNumber[int64]))...)
}

// Append 追加元素
func (i *Int8Array) Append(is ...int64) {
	if len(is) == 0 {
		return
	}

	slice := i.ToSlice()
	slice = append(slice, is...)
	i.Set(slice)
}

// AppendIfNotExist 追加元素，如果元素不存在
// 返回是否追加成功
func (i *Int8Array) AppendIfNotExist(is ...int64) bool {
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
		i.Set(slice)
	}
	return flag
}

func (i *Int8Array) AppendStrs(is ...string) {
	i.Append(xarray.MapBy(is, xparse.WarpToMap(xparse.ToNumber[int64]))...)
}

func (i *Int8Array) NilToNew(is ...int64) *Int8Array {
	if i == nil {
		return NewInt8Array(is...)
	}
	return i
}

func (i *Int8Array) Contain(val int64) bool {
	return xarray.Contain(i.ToSlice(), val)
}
func (i *Int8Array) ContainBy(predicate func(item int64) bool) bool {
	return xarray.ContainBy(i.ToSlice(), predicate)
}

// Remove 移除元素
func (i *Int8Array) Remove(val int64) bool {
	if i == nil {
		return false
	}
	slice := xarray.Delete(i.ToSlice(), val)
	err := i.Set(slice)
	if err != nil {
		xlog.Errorf(nil, "Int8Array.Remove error: %v", err)
		return false
	}
	return true
}

func RemoveInt8Array(i *Int8Array, is ...int64) bool {
	if i == nil || len(is) == 0 {
		return true
	}
	slice := i.ToSlice()
	for _, x := range is {
		slice = xarray.Delete(slice, x)
	}
	if err := i.Set(slice); err != nil {
		xlog.Errorf(nil, "RemoveInt8Array error: %v", err)
		return false
	}
	return true
}

func AppendInt8Array(i *Int8Array, is ...int64) {
	if i == nil {
		i = NewInt8Array(is...)
	}
	i.Append(is...)
}

func AppendInt8ArrayIfNotExist(i *Int8Array, is ...int64) bool {
	if i == nil {
		i = NewInt8Array(is...)
		return true
	}

	return i.AppendIfNotExist(is...)
}

// MarshalJSON 实现 json.Marshaler 接口
func (i *Int8Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ToStrings())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Int8Array) UnmarshalJSON(data []byte) error {
	var strArray []string
	if err := json.Unmarshal(data, &strArray); err != nil {
		return err
	}
	ints := xarray.MapBy(strArray, xparse.WarpToMap(xparse.ToNumber[int64]))
	return i.Set(ints)
}

// GormDataType 实现 gorm.DataType 接口 适配类型 `gorm:"type:int8[]"`
func (*Int8Array) GormDataType() string {
	return "int8[]"
}
