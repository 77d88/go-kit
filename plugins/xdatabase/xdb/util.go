package xdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/basic/xreflect"
	"gorm.io/gorm"
)

type Params map[string]interface{}

// IsNotFound 检查给定的错误是否为 gorm.ErrRecordNotFound 错误
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

func WarpLike(str string) string {
	if str != "" {
		return "%" + str + "%"
	}
	return ""
}

func WarpLikeRight(str string) string {
	if str != "" {
		return str + "%"
	}
	return ""
}
func WarpLikeLeft(str string) string {
	if str != "" {
		return "%" + str
	}
	return ""
}

func NextId() int64 {
	return xid.NextId()
}

// FindIds 根据字段名获取值 默认不要 <= 0
func FindIds[T any](value []T, field string, union bool) []int64 {
	ids := make([]int64, 0)
	if len(value) == 0 {
		return ids
	}
	for _, a := range value {
		val := xreflect.GetFieldVal(a, field)
		if val == nil {
			continue
		}
		switch c := val.(type) {
		case int64:
			if c > 0 {
				ids = append(ids, c)
			}
		case *Int8Array:
			if !c.IsEmpty() {
				ids = append(ids, c.ToSlice()...)
			}
		case Int8Array:
			if !c.IsEmpty() {
				ids = append(ids, c.ToSlice()...)
			}
		}
	}
	if union {
		return xarray.Union(ids)
	}
	return ids
}

// SortByIds 根据指定id排序
func SortByIds[T KeyModel](list []T, ids []int64) []T {
	ts := make([]T, 0)
	for _, c := range ids {
		for _, l := range list {
			id := l.GetID()
			if id == c {
				ts = append(ts, l)
				break
			}
		}
	}
	return ts
}

// FindLinksSet 一般用于关联设值
func FindLinksSet[T any, C comparable](slice []T, v C, get func(T) C, set func(T)) bool {
	x, ok := xarray.Find(slice, func(index int, item T) bool {
		return get(item) == v
	})
	if ok {
		set(x)
	}

	return ok
}

func FastDsn(host string, port int, user, password, dbName string) string {
	//host=pgsql port=5432 user=postgres password=123456 dbname=hospital sslmode=disable
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 5432
	}
	if password == "" {
		password = "postgres"
	}
	if dbName == "" {
		dbName = "postgres"
	}
	if user == "" {
		user = "postgres"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

func Param(name string, val any) sql.NamedArg {
	return sql.Named(name, val)
}
