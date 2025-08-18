package xdb

import (
	"github.com/77d88/go-kit/basic/xarray"
	"gorm.io/gorm"
)

type FindResult[T any] struct {
	Total int64
	List  []T
	Error error
}

func (t *FindResult[T]) IsEmpty() bool {
	return t.Error != nil || len(t.List) == 0
}

// FindPage 分页查询
func FindPage[T any](db *gorm.DB, page Pager, count bool) FindResult[T] {
	var pageResult FindResult[T]
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

// FindLinks 查询关联Id集合 fields 支持 int8Array 和 int64 其余不支持 默认去重查询 不保证排序
// T 可以为集合
// U 为关联对象 实际数据库查询对象
// fields 为关联字段
// 示例：xdb.DB().FindLinks[models.Order,models.User](&orders,"UserID","SendId")
// 上述实例 查询 订单中的 userID 和 sendId 忽略掉重复项 返回 用户集合
func FindLinks[T any, U any](db *gorm.DB, source []T, fields ...string) *FindResult[U] {
	// source 是否为集合
	ids := make([]int64, 0)
	for _, field := range fields {
		findIds := FindIds(source, field, false)
		ids = append(ids, findIds...)
	}
	ids = xarray.Union(ids)

	if len(ids) == 0 {
		return &FindResult[U]{}
	}
	var u []U
	tx := db.Where("id in (?)", ids).Find(&u)
	return &FindResult[U]{
		Error: tx.Error,
		List:  u,
	}
}
