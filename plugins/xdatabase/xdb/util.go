package xdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

type PageResult[T any] struct {
	Total int64
	List  []T
	Error error
}

func FindPage[T any](db *gorm.DB, page Pager, count bool) PageResult[T] {
	var pageResult PageResult[T]
	offset, limit := page.Limit()
	if count {
		if result := db.Count(&pageResult.Total); result.Error != nil {
			pageResult.Error = result.Error
			return pageResult
		}
		if pageResult.Total <= int64(offset) {
			return pageResult
		}
	}
	find := db.Offset(offset).Limit(limit).Find(&pageResult.List)
	pageResult.Error = find.Error
	return pageResult
}

func Session(db *gorm.DB, session ...*gorm.Session) *gorm.DB {
	return db.Session(xarray.FirstOrDefault(session, &gorm.Session{}))
}

func XWhere(db *gorm.DB, condition bool, query string, args ...interface{}) *gorm.DB {
	if query == "" {
		return db
	}
	if !condition {
		return db
	}
	return db.Where(query, args...)
}

type SaveMapResult struct {
	Error error
	RowId int64
}

func SaveMap[T any](db *gorm.DB, obj interface{}, mapping ...interface{}) *SaveMapResult {
	m := toSqlMap(obj, mapping...)
	mdb := db.Model(new(T))
	var id int64
	var r SaveMapResult
	// 获取出ID 单独处理
	for k, v := range m {
		if strings.ToLower(k) == "id" {
			delete(m, k)
			if i, ok := v.(int64); ok {
				id = i
			}
		}
	}
	if id > 0 {
		result := mdb.Where("id = ?", id).Updates(m)
		r.RowId = id
		r.Error = result.Error
	} else {
		saveId := NextId()
		m["id"] = saveId
		m["created_time"] = time.Now()
		m["updated_time"] = time.Now()
		m["deleted_time"] = nil
		result := mdb.Create(m)
		r.RowId = saveId
		r.Error = result.Error
	}
	return &r
}
