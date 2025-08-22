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

	// 创建通道用于接收计数结果
	countChan := make(chan struct {
		total int64
		err   error
	}, 1)

	// 创建通道用于接收查询结果
	dataChan := make(chan struct {
		list []T
		err  error
	}, 1)

	// 异步执行 count 查询
	if count {
		go func() {
			var total int64
			var err error
			if result := db.Session(&gorm.Session{}).Model(new(T)).Count(&total); result.Error != nil {
				err = result.Error
			}
			countChan <- struct {
				total int64
				err   error
			}{total: total, err: err}
		}()
	}

	// 异步执行数据查询
	go func() {
		var list []T
		var err error
		result := db.Session(&gorm.Session{}).Model(new(T)).Offset(offset).Limit(limit).Find(&list)
		if result.Error != nil {
			err = result.Error
		}
		dataChan <- struct {
			list []T
			err  error
		}{list: list, err: err}
	}()

	// 等待并处理数据查询结果
	dataResult := <-dataChan
	pageResult.List = dataResult.list
	pageResult.Error = dataResult.err

	// 如果数据查询出错，直接返回
	if pageResult.Error != nil {
		return pageResult
	}

	// 等待并处理 count 结果
	if count {
		countResult := <-countChan
		if countResult.err != nil {
			pageResult.Error = countResult.err
			return pageResult
		}
		pageResult.Total = countResult.total

		// 如果总数为0或偏移量超过总数，清空列表
		if pageResult.Total <= int64(offset) {
			pageResult.List = []T{}
		}
	}

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
