package xpg

import (
	"errors"
	"fmt"
	"reflect"

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
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		clone := r.Clone()
		clone.Error = errors.Join(r.Error, fmt.Errorf("result must be a pointer"))
		return clone
	}

	switch i.(type) {
	case *int, *int64, *float64, *string, *bool:
		covert, ok := mapFirstCovert(r.MapResult, i)
		if ok {
			// 使用反射正确设置值
			reflect.ValueOf(i).Elem().Set(reflect.ValueOf(covert))
			r.Result = covert
			return r
		} else {
			r.Error = fmt.Errorf("covert to %v error", reflect.TypeOf(i))
			return r
		}
	}
	n := r.Clone()
	n.Error = Scan(n.MapResult, i)
	if n.Error == nil {
		n.Result = i
	}
	return n
}

func (r *Result) Clone() *Result {
	return &Result{
		Args:      r.Args,
		Error:     r.Error,
		MapResult: r.MapResult,
		Result:    r.Result,
		RowId:     r.RowId,
		Rows:      r.Rows,
		Sql:       r.Sql,
		Total:     r.Total,
	}

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

// Salvage 捞取关联数据 简单的 如果target不是model ins 需要提前设置 tablename和捞取字段
func (r *Result) Salvage(ins *Inst, target interface{}, fields ...string) *Result {
	ids := r.SalvageId(fields...)
	if len(ids) == 0 {
		return r
	}
	return ins.Where("id = any(?)", ids).Find(target)
}

// SalvageId 捞取ID
func (r *Result) SalvageId(fields ...string) []int64 {
	// 捞取ID
	var ids []int64
	unique := make(map[int64]struct{})
	for _, m := range r.MapResult {
		for _, field := range fields {
			if id, ok := m[field]; ok {
				if i, ok := id.(int64); ok {
					if _, exists := unique[i]; !exists {
						unique[i] = struct{}{}
						ids = append(ids, i)
					}
				}
			}
		}
	}
	return ids
}
