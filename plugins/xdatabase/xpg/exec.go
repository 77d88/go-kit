package xpg

import (
	"errors"
	"strings"
	"time"

	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/basic/xtime"
	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (i *Inst) Exec(sql string, args ...interface{}) (res *Result) {
	if i.debug || i.config.Logger {
		inv := xtime.NewTimeInterval()
		defer func() {
			xlog.Debugf(i.ctx, "Exec[%dms】 RowsAffected:%d sql:%s args:%v", inv.IntervalMs(), res.Rows, sql, args)
		}()
	}
	var conn pgconn.CommandTag
	var err error
	if i.tx != nil {
		conn, err = i.tx.Exec(i.ctx, sql, args...)
	} else {
		conn, err = i.pool.Exec(i.ctx, sql, args...)
	}
	if err != nil {
		xlog.Errorf(i.ctx, "handlerExec error: %v", err)
		return &Result{Error: err}
	}
	return &Result{Error: err, Rows: conn.RowsAffected()}
}

func (i *Inst) Update(field string, value any) *Result {
	return i.Updates(map[string]interface{}{field: value})
}

func (i *Inst) Updates(m map[string]interface{}) *Result {
	if i.tableName == "" {
		return &Result{Error: errors.New("tableName is empty")}
	}
	clause, i2, err := i.extractWhereClause()
	if err != nil {
		return &Result{Error: err}
	}
	if clause == "" {
		return &Result{Error: errors.New("where clause is empty")}
	}
	where := sq.Update(i.tableName).Where(clause, i2...)
	for field, value := range m {
		where = where.Set(field, value)
	}
	// 更新时间
	where = where.Set("updated_time", time.Now())

	sql, i3, err := where.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return &Result{Error: err}
	}
	return i.Exec(sql, i3...)
}

func (i *Inst) extractWhereClause() (string, []interface{}, error) {
	// 获取 select 语句的 where 部分
	selectSQL, selectArgs, err := i.cond.Columns("1").PlaceholderFormat(sq.Question).ToSql()
	if err != nil {
		return "", nil, err
	}
	// 实现 SQL 解析逻辑来提取 WHERE 子句
	// 这是一个简化的示例
	whereIndex := strings.Index(strings.ToUpper(selectSQL), " WHERE ")
	if whereIndex == -1 {
		return "", nil, nil
	}
	return selectSQL[whereIndex+7:], selectArgs, nil // 7 是 " WHERE " 的长度
}

func (i *Inst) Create(c interface{}) *Result {
	if i.tableName == "" {
		return &Result{Error: errors.New("tableName is empty")}
	}
	dbObj, err := extractDBObj(c)
	if err != nil {
		return &Result{Error: err}
	}
	// 默认字段设定
	dbObj["id"] = xid.NextIdStr()
	dbObj["created_time"] = time.Now()
	dbObj["updated_time"] = time.Now()

	// 去除map里面的空字段
	keys := make([]string, 0, len(dbObj))
	values := make([]any, 0, len(dbObj))
	for k, v := range dbObj {
		if xcore.IsZero(v) { //忽略空字段
			continue
		} else {
			keys = append(keys, k)
			values = append(values, v)
		}
	}

	insert := sq.Insert(i.tableName).Columns(keys...)
	insert = insert.Values(values...)
	sql, i3, err := insert.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return &Result{Error: err}
	}
	return i.Exec(sql, i3...)
}

// Save 保存 如果id存在则更新，不存在则创建
func (i *Inst) Save(obj interface{}, before ...func(m map[string]interface{})) *Result {
	ic := i.Copy()
	if m, ok := obj.(Model); ok && i.tableName == "" {
		ic.tableName = m.TableName()
	}
	dbObj, err := extractDBObj(obj)
	if err != nil {
		return &Result{Error: err}
	}
	//忽略空字段
	for k, v := range dbObj {
		if xcore.IsZero(v) {
			delete(dbObj, k)
		}
	}
	if id, ok := dbObj["id"]; ok {
		delete(dbObj, "id") // 忽略id字段
		for _, f := range before {
			f(dbObj)
		}
		return ic.Where("id=?", id).Updates(dbObj)
	} else {
		for _, f := range before {
			f(dbObj)
		}
		return ic.Create(obj)
	}
}
