package xdb

import (
	"gorm.io/gorm"
)

var emptyResult = &Result{}

type Result struct {
	Error        error
	RowsAffected int64
	RowId        int64
}

func (r *Result) GetError() error {
	if r == nil {
		return nil
	}
	return r.Error
}

// IsNotFound 是否是没有数据 调用First的时候在使用这个
func (r *Result) IsNotFound() bool {
	return IsNotFound(r.Error)
}

func warpResult(db *gorm.DB) *Result {
	return &Result{
		Error:        db.Error,
		RowsAffected: db.RowsAffected,
	}
}
