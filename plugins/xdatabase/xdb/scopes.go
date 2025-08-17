package xdb

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xreflect"
	"gorm.io/gorm"
)

type Pager interface {
	Limit() (offset, limit int)
}

type SC_Opts int

const (
	SC_IgnoreNull SC_Opts = iota
	SC_IgnoreZero
	SC_OrderIdDesc
)

var sc_empty = func(*gorm.Statement) {}


// SC_ZeroWhere 忽略空值的条件查询
// 示例：
//
//	SC_ZeroIgnore(0, func(db *gorm.db, r int){
//			db.Where("id = ?", r)
//		}) 这样是不会添加 id = 0 这个条件
func SC_ZeroWhere(value any, query interface{}, args ...interface{}) func(*gorm.Statement) {
	if xcore.IsZero(value) {
		return sc_empty
	}
	return func(s *gorm.Statement) {
		s.Where(query, args...)
	}
}

// SC_NilWhere 忽略 nil 的条件查询 和 SC_ZeroWhere 差不多只是条件值可以为零值
func SC_NilWhere(value any, query interface{}, args ...interface{}) func(*gorm.Statement) {
	if xreflect.IsNil(value) {
		return sc_empty
	}
	return func(s *gorm.Statement) {
		s.Where(query, args...)
	}
}

func SC_Page(page, size int) func(*gorm.Statement) {
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}
	return func(s *gorm.Statement) {
		s.Offset((page - 1) * size).Limit(size)
	}
}

func SC_Ilike(field, page string) func(*gorm.Statement) {
	return func(s *gorm.Statement) {
		s.Where(field+" ILIKE ?", "%"+page+"%")
	}
}

func SC_ID(id int64) func(*gorm.Statement) {
	return func(s *gorm.Statement) {
		s.Where("id =  ?", id)
	}
}
