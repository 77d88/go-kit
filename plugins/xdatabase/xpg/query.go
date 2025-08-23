package xpg

import (
	"context"
	"fmt"
	"reflect"

	"github.com/77d88/go-kit/basic/xtime"
	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/lann/builder"
)

func (i *Inst) First(result interface{}) *Result {
	// 检查result类型来决定是查询单个还是多个记录
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Ptr {
		return &Result{
			Error: fmt.Errorf("result must be a pointer"),
		}
	}
	// 如果是切片指针，则查询多个记录
	if t.Elem().Kind() == reflect.Slice {
		return &Result{
			Error: fmt.Errorf("result must be a pointer to a slice of structs"), // 必须是struct 类型
		}
	}
	// 只查询一个
	find := i.Find(result)
	if find.Rows == 0 {
		return &Result{
			Error: pgx.ErrNoRows,
		}
	}
	return find
}

func (i *Inst) Find(result interface{}) *Result {
	// 检查result类型来决定是查询单个还是多个记录
	ic := i
	t := reflect.TypeOf(result)
	if t.Kind() != reflect.Ptr {
		return &Result{
			Error: fmt.Errorf("result must be a pointer"),
		}
	}

	// 如果是切片指针，则查询多个记录
	if t.Elem().Kind() != reflect.Slice {
		// 检查是否是model
		if ic.tableName == "" || ic.model == nil {
			if m, ok := result.(Model); ok {
				ic = ic.Model(m)
			}
		}

		ic.cond = ic.cond.Limit(1)
	} else {
		// 切片的话使用类型初始化一个进行 mode
		if ic.tableName == "" || ic.model == nil {
			structType := t.Elem().Elem()
			if structType.Kind() == reflect.Ptr {
				structType = structType.Elem()
			}
			a := reflect.New(structType).Interface()
			if m, ok := a.(Model); ok {
				ic = ic.Model(m)
			}
		}
	}

	rs, err := ic.Query()
	if err != nil {
		return &Result{
			Error: err,
		}
	}
	if len(rs) == 0 {
		return &Result{}
	}
	err = Scan(rs, result)
	return &Result{
		Rows:   int64(len(rs)),
		Error:  err,
		Result: result,
	}
}

func (i *Inst) Count(rc ...*int64) *Result {
	ic := i.Select(`COUNT(1) as count`)
	ic.cond = builder.Delete(ic.cond, "OrderByParts").(sq.SelectBuilder)
	rs, err := ic.Query()
	if err != nil {
		return &Result{
			Error: err,
		}
	}

	m := rs[0]
	if m == nil {
		return &Result{
			Error: fmt.Errorf("count error"),
		}
	}
	if m["count"] == nil {
		return &Result{
			Error: fmt.Errorf("count error"),
		}
	}
	count, ok := m["count"].(int64)
	if !ok {
		return &Result{
			Error: fmt.Errorf("count error"),
		}
	}
	if len(rc) > 0 {
		*rc[0] = count
	}
	return &Result{
		Rows:   1,
		Error:  nil,
		Total:  count,
		Result: count,
	}
}

// Query 查询结果 并统一处理为 map
func (i *Inst) Query() ([]map[string]any, error) {
	maps := make([]map[string]any, 0)
	sql, args, err := i.cond.From(i.tableName).Columns(i.selectFields...).Where("deleted_time is null").ToSql() // 都是软删除
	if err != nil {
		return maps, err
	}
	if i.debug || i.config.Logger {
		inv := xtime.NewTimeInterval()
		defer func() {
			xlog.Debugf(i.ctx, "Query[%dms rows[%d]] sql: %s, args: %v", inv.IntervalMs(), len(maps), sql, args)
		}()
	}
	rows, err := i.pool.Query(i.ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return maps, err
	}

	for rows.Next() {
		toMap, err := pgx.RowToMap(rows)
		if err != nil {
			xlog.Errorf(context.Background(), "handlerQueryRowToMap error: %v", err)
			return maps, err
		}
		maps = append(maps, toMap)
	}
	return maps, nil
}

// FindPage 分页查询
func (i *Inst) FindPage(list interface{}, page Pager, count bool) *Result {
	offset, limit := page.Limit()
	c1 := make(chan *Result)
	if count {
		go func(inst *Inst) {
			c1 <- inst.Count()
		}(i.Copy())
	}
	inst := i.Limit(limit)
	if offset > 0 {
		inst = inst.Offset(offset)
	}
	find := inst.Find(list)
	if count {
		r := <-c1
		find.Total = r.Total
	}
	return find
}
