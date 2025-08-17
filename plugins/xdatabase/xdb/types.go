package xdb

//	type SearchRequest struct {
//		Val string `form:"val" json:"val"` // 搜索关键字
//		Ids string `form:"ids" json:"ids"` // 默认数据ID 1,2,3 形式
//	}
//
//	func (r SearchRequest) ToIntIds() []int64 {
//		ids, err := xstr.SplitTo(r.Ids, ",", xparse.ToNumber[int64])
//		if err != nil {
//			return make([]int64, 0)
//		}
//		return xarray.Union(ids)
//	}
//
//	type PageSearch struct {
//		Page PageSearch `form:"page" json:"page"`
//	}
//
//	func (a PageSearch) Limit() (offset, limit int) {
//		return a.Page.Limit()
//	}
//
//	func (a PageSearch) IsNotCounted() bool {
//		return a.Page.IsNotCounted()
//	}
//
//	type Pager interface {
//		Limit() (offset, limit int)
//		IsNotCounted() bool
//	}
type PageSearch struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p PageSearch) Limit() (offset, limit int) {
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

//func (p PageSearch) IsNotCounted() bool {
//	return p.NotCounted
//}
