package xpg

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

type Result struct {
	Error  error
	Rows   int64
	RowId  int64
	Total  int64       // 统计内置
	Result interface{} // 如果有scan的结果也放在这里可以通过
}

func (r *Result) IsNotFound() bool {
	return errors.Is(r.Error, pgx.ErrNoRows)
}

func (r *Result) Decon() (interface{}, error) {
	return r.Result, r.Error
}
