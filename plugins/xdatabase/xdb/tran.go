package xdb

import (
	"context"
	"database/sql"

	"github.com/77d88/go-kit/basic/xctx"
	"gorm.io/gorm"
)

const CtxTransactionKey = "__CTX_TRANSACTION_KEY__"

type TranOpt struct {
	SqlOpts      *sql.TxOptions
	DataBaseName string // 数据源名称
}

func (d *DB) Tran(fc func(tx *DB) error, opts ...*sql.TxOptions) (err error) {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		return fc(d.wrap(tx))
	}, opts...)

}

func (d *DB) Begin(opts ...*sql.TxOptions) *DB {
	return d.wrap(d.DB.Begin(opts...))
}

func (d *DB) Commit() *DB {
	return d.wrap(d.DB.Commit())
}

func (d *DB) SavePoint(name string) *DB {
	return d.wrap(d.DB.SavePoint(name))
}

func (d *DB) RollbackTo(name string) *DB {
	return d.wrap(d.DB.RollbackTo(name))
}

func GetCtxTran(c context.Context) *DB {
	get := c.Value(CtxTransactionKey) // 有事务优先获取事务
	if get != nil {
		return get.(*DB)
	}
	return nil
}

// BeginWithCtx 在当前上下文中开始事务
func BeginWithCtx(c context.Context, opts ...*sql.TxOptions) *DB {
	if tx := GetCtxTran(c); tx != nil {
		return tx // 如果已有事务，返回现有事务
	}
	begin := Ctx(c).Begin(opts...)
	xctx.SetVal(c, CtxTransactionKey, begin)
	return begin
}

func BeginWithCtxAndDbName(c context.Context, name string, opts ...*sql.TxOptions) *DB {
	if tx := GetCtxTran(c); tx != nil {
		return tx // 如果已有事务，返回现有事务
	}
	begin := Ctx(c, name).Begin(opts...)
	xctx.SetVal(c, CtxTransactionKey, begin)
	return begin
}
