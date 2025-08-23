package xtype

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Int64Array []int64

func (i Int64Array) ToSlice() []int64 {
	return i
}
func (i Int64Array) IsEmpty() bool {
	return len(i) == 0
}

func (i Int64Array) ToStrings() []string {
	// 初始化一个空的字符串数组
	var goArray = make([]string, 0, len(i))
	for i, elem := range i {
		goArray[i] = strconv.FormatInt(elem, 10)
	}
	return goArray
}

// MarshalJSON 实现 json.Marshaler 接口
func (i Int64Array) MarshalJSON() ([]byte, error) {
	// 将 int64 转换为字符串数组
	strs := make([]string, len(i))
	for idx, val := range i {
		strs[idx] = strconv.FormatInt(val, 10)
	}
	return json.Marshal(strs)
}

// UnmarshalJSON 从字符串数组反序列化为 int64 数组
func (i *Int64Array) UnmarshalJSON(data []byte) error {
	if i == nil {
		return errors.New("nil pointer")
	}
	var strArray []string
	if err := json.Unmarshal(data, &strArray); err != nil {
		return err
	}

	// 将字符串数组转换为 int64 数组
	ints := make([]int64, len(strArray))
	for idx, str := range strArray {
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		ints[idx] = val
	}

	*i = ints
	return nil
}

func (i Int64Array) Contain(cp int64) bool {
	for _, v := range i {
		if v == cp {
			return true
		}
	}
	return false
}
