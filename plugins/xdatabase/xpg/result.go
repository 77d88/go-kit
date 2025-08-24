package xpg

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

type Result struct {
	Error     error            // 错误
	Rows      int64            // 查询结果行数 or 影响行数
	RowId     int64            // 插入时返回的id
	Total     int64            // 统计内置
	Result    interface{}      // 如果有scan的结果也放在这里可以通过
	MapResult []map[string]any // 原始查询结果 转为了map
	Sql       string           // 原始sql
	Args      []any            // 参数
}

// AddError 添加错误
func (r *Result) AddError(err error) *Result {
	r.Error = errors.Join(r.Error, err)
	return r
}

// IsNotFound 判断是否是找不到 First的时候才会有
func (r *Result) IsNotFound() bool {
	return errors.Is(r.Error, pgx.ErrNoRows)
}

// Decon 解包 为 返回结果和错误
func (r *Result) Decon() (interface{}, error) {
	return r.Result, r.Error
}

// Scan 扫描结果将原始map扫描到传入的结构体中 并保存到Result.Result中
func (r *Result) Scan(i interface{}) *Result {
	r.Error = Scan(r.MapResult, i)
	if r.Error == nil {
		r.Result = i
	}
	return r
}

// Get 获取原始map的某个字段
// 索引从0开始 i 为行数 key 为字段名
func (r *Result) Get(i int, key string) (interface{}, bool) {
	m := r.MapResult[i]
	if m == nil {
		return nil, false
	}
	x, exist := m[key]
	return x, exist
}
