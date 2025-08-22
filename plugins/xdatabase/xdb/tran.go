package xdb

//
//const CtxTransactionKey = "__CTX_TRANSACTION_KEY__"
//
//type TranOpt struct {
//	SqlOpts      *sql.TxOptions
//	DataBaseName string // 数据源名称
//}
//
//func (d *db) Tran(fc func(tx *db) error, opts ...*sql.TxOptions) (err error) {
//	return d.db.Transaction(func(tx *gorm.db) error {
//		return fc(d.wrap(tx))
//	}, opts...)
//
//}
//
//func (d *db) Begin(opts ...*sql.TxOptions) *db {
//	return d.wrap(d.db.Begin(opts...))
//}
//
//func (d *db) Commit() *db {
//	return d.wrap(d.db.Commit())
//}
//
//func (d *db) SavePoint(name string) *db {
//	return d.wrap(d.db.SavePoint(name))
//}
//
//func (d *db) RollbackTo(name string) *db {
//	return d.wrap(d.db.RollbackTo(name))
//}
//
//func GetCtxTran(c context.Context) *db {
//	get := c.Value(CtxTransactionKey) // 有事务优先获取事务
//	if get != nil {
//		return get.(*db)
//	}
//	return nil
//}
//
//// BeginWithCtx 在当前上下文中开始事务
//func BeginWithCtx(c context.Context, opts ...*sql.TxOptions) *db {
//	if tx := GetCtxTran(c); tx != nil {
//		return tx // 如果已有事务，返回现有事务
//	}
//	begin := C(c).Begin(opts...)
//	xctx.SetVal(c, CtxTransactionKey, begin)
//	return begin
//}
//
//func BeginWithCtxAndDbName(c context.Context, name string, opts ...*sql.TxOptions) *db {
//	if tx := GetCtxTran(c); tx != nil {
//		return tx // 如果已有事务，返回现有事务
//	}
//	begin := C(c, name).Begin(opts...)
//	xctx.SetVal(c, CtxTransactionKey, begin)
//	return begin
//}

//func (db *DB) Begin() *DB {
//	newDB := db.clone()
//	tx, err := newDB.pool.Begin(newDB.ctx)
//	if err != nil {
//		newDB.error = err
//		return newDB
//	}
//	newDB.tx = tx
//	return newDB
//}
//
//func (db *DB) Commit() error {
//	if db.tx != nil {
//		return db.tx.Commit(db.ctx)
//	}
//	return nil
//}
//
//func (db *DB) Rollback() error {
//	if db.tx != nil {
//		return db.tx.Rollback(db.ctx)
//	}
//	return nil
//}
