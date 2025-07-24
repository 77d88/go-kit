package xdb

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

const CtxTransactionKey = "__CTX_TRANSACTION_KEY__"

type TranOpt struct {
	SqlOpts      *sql.TxOptions
	DataBaseName string // 数据源名称
}

func (d *DataSource) Tran(fc func(tx *DataSource) error, opts ...*sql.TxOptions) (err error) {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		return fc(wrap(d.Context, tx))
	}, opts...)

}

func CtxTran(c context.Context, f func(d *DataSource) error, opts ...*TranOpt) error {
	var opt *TranOpt
	if len(opts) > 0 {
		opt = opts[0]
	}
	name := dbStr
	var sqlOpts = make([]*sql.TxOptions, 0)
	if opt != nil {
		name = opt.DataBaseName
		sqlOpts = append(sqlOpts, opt.SqlOpts)
	}

	return Ctx(c, name).Tran(func(tx *DataSource) error {
		tx.Context = context.WithValue(c, CtxTransactionKey, tx)
		return f(tx)
	}, sqlOpts...)
}
