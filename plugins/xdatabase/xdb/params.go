package xdb

import (
	"github.com/77d88/go-kit/basic/xcore"
	"gorm.io/gorm"
)

type ParamBuilder[T any] struct {
	params map[string]interface{}
}

// ExtZero 添加参数，如果参数为零值则不添加
func (p *ParamBuilder[T]) ExtZero(query string, value interface{}) *ParamBuilder[T] {
	if xcore.IsZero(value) {
		return p
	}
	return p.Add(query, value)
}

func (p *ParamBuilder[T]) Add(query string, value interface{}) *ParamBuilder[T] {
	p.params[query] = value
	return p
}

func (p *ParamBuilder[T]) Eq(query string, value interface{}) *ParamBuilder[T] {
	return p.Add(query+" = ?", value)
}

func (p *ParamBuilder[T]) EqNonZero(query string, value interface{}) *ParamBuilder[T] {
	return p.ExtZero(query+" = ?", value)
}

func (p *ParamBuilder[T]) ILike(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" ILIKE ?", WarpLike(value))
}
func (p *ParamBuilder[T]) ILikeLeft(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" ILIKE ?", WarpLikeLeft(value))
}
func (p *ParamBuilder[T]) ILikeRight(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" ILIKE ?", WarpLikeRight(value))
}
func (p *ParamBuilder[T]) Like(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" LIKE ?", WarpLike(value))
}
func (p *ParamBuilder[T]) LikeLeft(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" LIKE ?", WarpLikeLeft(value))
}
func (p *ParamBuilder[T]) LikeRight(query string, value string) *ParamBuilder[T] {
	return p.ExtZero(query+" LIKE ?", WarpLikeRight(value))
}


func (p *ParamBuilder[T]) Build(db *gorm.DB) *gorm.DB {
	scops := make([]func(*gorm.DB) *gorm.DB, 0, len(p.params))
	for query, value := range p.params {
		scops = append(scops, func(gdb *gorm.DB) *gorm.DB {
			return gdb.Where(query, value)
		})
	}
	return db.Scopes(scops...)
}

// BuildC 这个是泛型兼容版本
func (p *ParamBuilder[T]) BuildC(db *gorm.DB) gorm.ChainInterface[T] {
	scops := make([]func(*gorm.Statement), 0, len(p.params))
	for query, value := range p.params {
		scops = append(scops, func(statement *gorm.Statement) {
			statement.Where(query, value)
		})
	}
	return gorm.G[T](db).Scopes(scops...)
}

func NewParams[T any]() *ParamBuilder[T] {
	return &ParamBuilder[T]{
		params: make(map[string]interface{}),
	}
}
