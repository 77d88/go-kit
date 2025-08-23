package xpg

type Pager interface {
	Limit() (offset, limit int)
}

// PageSearch 分页查询的
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
