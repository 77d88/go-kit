package xdb

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/basic/xstr"
)

type SearchRequest struct {
	Val string `form:"val" json:"val"` // 搜索关键字
	Ids string `form:"ids" json:"ids"` // 默认数据ID 1,2,3 形式
}

func (r SearchRequest) ToIntIds() []int64 {
	ids, err := xstr.SplitTo(r.Ids, ",", xparse.ToNumber[int64])
	if err != nil {
		return make([]int64, 0)
	}
	return xarray.Union(ids)
}

type PageSearch struct {
	Page PageRequestImpl `form:"page" json:"page"`
}

func (a PageSearch) Limit() (offset, limit int) {
	return a.Page.Limit()
}

func (a PageSearch) IsNotCounted() bool {
	return a.Page.IsNotCounted()
}

type PageRequest interface {
	Limit() (offset, limit int)
	IsNotCounted() bool
}

type PageRequestImpl struct {
	Page,
	Size int
	NotCounted bool
}

func (p PageRequestImpl) Limit() (offset, limit int) {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size > 1000 || p.Size < 1 { // 每页最大1000
		p.Size = 20
	}
	offset = (p.Page - 1) * p.Size
	limit = p.Size
	return
}
func (p PageRequestImpl) IsNotCounted() bool {
	return p.NotCounted
}
