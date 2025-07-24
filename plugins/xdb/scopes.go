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

var sc_empty = func(db *gorm.DB) *gorm.DB {
	return db
}

// SC_ZeroIgnore 忽略空值的条件查询
// 示例：
//
//	SC_ZeroIgnore(0, func(db *gorm.DB, r int){
//			db.Where("id = ?", r)
//		}) 这样是不会添加 id = 0 这个条件
func SC_ZeroIgnore[T any](value T, f func(db *gorm.DB, v T) *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if xcore.IsZero(value) {
		return sc_empty
	}
	return func(db *gorm.DB) *gorm.DB {
		return f(db, value)
	}
}

// SC_NullIgnore 忽略 nil 的条件查询 和 SC_ZeroIgnore 差不多只是条件值可以为零值
// 示例：SC_NullIgnore(0, func(db *gorm.DB, r int){ db.Where("id = ?", r)  }) 这样会添加 id = 0 这个条件
// 示例：SC_NullIgnore(nil, func(db *gorm.DB, r int){ db.Where("id = ?", r)  }) 这样不会添加 id = 0 这个条件
func SC_NullIgnore[T any](value T, f func(db *gorm.DB, v T) *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if xreflect.IsNil(value) {
		return sc_empty
	}
	return func(db *gorm.DB) *gorm.DB {

		return f(db, value)
	}
}
