package xpg

import (
	"fmt"
	"reflect"

	"github.com/77d88/go-kit/basic/xtime"
	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/lann/builder"
)

// Raw 执行sql 并统一处理为 map
func (i *Inst) Raw(sql string, args ...interface{}) (re *Result) {
	if i.debug || i.config.Logger {
		inv := xtime.NewTimeInterval()
		defer func() {
			xlog.Debugf(i.ctx, "Exec[%dms rows[%d] tx:[%s]] sql: %s, args: %v", inv.IntervalMs(), re.Rows, sql, args, i.savepointCounter)
		}()
	}
	var rows pgx.Rows
	var queryErr error
	if i.tx != nil {
		rows, queryErr = i.tx.Query(i.ctx, sql, args...)
	} else {
		rows, queryErr = i.pool.Query(i.ctx, sql, args...)
	}
	if queryErr != nil {
		return &Result{Error: queryErr, Sql: sql, Args: args}
	}
	maps, err := RowsToMaps(i.ctx, rows)
	re = &Result{MapResult: maps, Error: err, Rows: int64(len(maps)), Sql: sql, Args: args}
	return
}

// Query 处理内置条件 并统一处理为 map
func (i *Inst) Query() *Result {
	sql, args, err := i.cond.From(i.tableName).Columns(i.selectFields...).Where("deleted_time is null").ToSql() // 都是软删除
	if err != nil {
		return &Result{Error: err}
	}
	return i.Raw(sql, args...)
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
	return ic.Query().Scan(result)
}
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

func (i *Inst) Count(rc ...*int64) *Result {
	ic := i
	if len(i.selectFields) == 0 {
		ic = i.Select(`COUNT(1)`)
	}
	ic.cond = builder.Delete(ic.cond, "OrderByParts").(sq.SelectBuilder)
	result := ic.Query()
	if result.Result != nil {
		return result
	}
	var count int64
	result = result.Scan(&count)
	if len(rc) > 0 {
		*rc[0] = count
	}
	result.Total = count
	return result
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
