package xdb2

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xtime"
	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type Inst struct {
	pool             *pgxpool.Pool
	cond             sq.SelectBuilder
	ctx              context.Context
	tx               pgx.Tx
	savepointCounter int
	spMu             sync.Mutex
	txMarkedBad      bool
}

func (i *Inst) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if i.tx != nil {
		return i.tx.Exec(i.ctx, sql, args...)
	}
	return i.pool.Exec(i.ctx, sql, args...)
}

// WithContext 设置上下文
func (i *Inst) WithContext(ctx context.Context) *Inst {
	return &Inst{
		pool:             i.pool,
		cond:             sq.Select().PlaceholderFormat(sq.Dollar),
		ctx:              ctx,
		tx:               i.tx,
		savepointCounter: i.savepointCounter,
	}
}
func (i *Inst) Table(table string) *Inst {
	i.cond = i.cond.From(table)
	return i
}

func (i *Inst) Select(fields ...string) *Inst {
	i.cond = i.cond.Columns(fields...)
	return i
}

// Where 添加 WHERE 条件
func (i *Inst) Where(query interface{}, args ...interface{}) *Inst {
	i.cond = i.cond.Where(query, args...)
	return i
}

func (i *Inst) Limit(limit int) *Inst {
	i.cond = i.cond.Limit(uint64(limit))
	return i
}

func (i *Inst) Offset(offset int) *Inst {
	i.cond = i.cond.Offset(uint64(offset))
	return i
}

func (i *Inst) Order(order string) *Inst {
	i.cond = i.cond.OrderBy(order)
	return i
}

func (i *Inst) Find(result interface{}) error {

	// 检查result类型来决定是查询单个还是多个记录
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}

	// 如果是切片指针，则查询多个记录
	if resultValue.Elem().Kind() != reflect.Slice {
		i.cond = i.cond.Limit(1)

	}
	sql, args, err := i.cond.ToSql()
	if err != nil {
		return err
	}
	inv := xtime.NewTimeInterval()
	query, err := i.pool.Query(i.ctx, sql, args...)
	xlog.Debugf(i.ctx, "sql【%d】ms: %s, args: %v", inv.IntervalMs(), sql, args)
	defer query.Close()
	if err != nil {
		return err
	}
	inv = xtime.NewTimeInterval()
	toMap := handlerQueryRowToMap(query)
	xlog.Debugf(i.ctx, "sql【%d】ms: %s, args: %v", inv.IntervalMs(), sql, args)

	xlog.Debugf(i.ctx, "values: %d %v", len(toMap), toMap)
	return nil
}
func (i *Inst) getNextSavepoint() string {
	i.spMu.Lock()
	defer i.spMu.Unlock()
	i.savepointCounter++
	return fmt.Sprintf("sp_%d_%d", time.Now().UnixNano(), i.savepointCounter)
}

// Transaction 处理事务，支持真正的嵌套事务
func (i *Inst) Transaction(fn func(*Inst) error) error {
	// 如果已经在事务中，使用savepoint实现嵌套事务
	if i.tx != nil {
		if i.txMarkedBad {
			return errors.New("transaction already marked bad")
		}

		savepointName := i.getNextSavepoint()

		_, err := i.tx.Exec(i.ctx, "SAVEPOINT "+pq.QuoteIdentifier(savepointName))
		if err != nil {
			// 如果无法创建savepoint，标记事务为bad
			i.txMarkedBad = true
			return err
		}

		defer func() {
			if p := recover(); p != nil {
				_, err := i.tx.Exec(i.ctx, "ROLLBACK TO SAVEPOINT "+pq.QuoteIdentifier(savepointName))
				if err != nil {
					i.txMarkedBad = true
					xlog.Errorf(i.ctx, "critical: rollback to savepoint failed: %v", err)
					return
				}
				panic(p)
			}
		}()

		err = fn(i)
		if err != nil {
			_, rollbackErr := i.tx.Exec(i.ctx, "ROLLBACK TO SAVEPOINT "+pq.QuoteIdentifier(savepointName))
			if rollbackErr != nil {
				i.txMarkedBad = true
				_ = i.tx.Rollback(i.ctx)
				return fmt.Errorf("critical: rollback to savepoint failed: %w", rollbackErr)
			}
			return err
		}

		_, releaseErr := i.tx.Exec(i.ctx, "RELEASE SAVEPOINT "+pq.QuoteIdentifier(savepointName))
		if releaseErr != nil {
			xlog.Warnf(i.ctx, "warning: release savepoint failed: %v", releaseErr)
		}
		return nil
	}

	// 开始新事务
	begin, err := i.pool.Begin(i.ctx)
	if err != nil {
		return err
	}

	txInst := &Inst{
		pool: i.pool,
		cond: sq.SelectBuilder{},
		ctx:  i.ctx,
		tx:   begin,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = begin.Rollback(i.ctx)
			panic(p)
		}
	}()

	// 执行业务逻辑
	err = fn(txInst)

	// 处理错误和提交/回滚
	if err != nil {
		_ = begin.Rollback(i.ctx)
		return err
	}

	// 提交事务
	return begin.Commit(i.ctx)
}

func handlerQueryRowToMap(rows pgx.Rows) []map[string]any {
	defer rows.Close()
	maps := make([]map[string]any, 0)
	for rows.Next() {
		toMap, err := pgx.RowToMap(rows)
		if err != nil {
			xlog.Errorf(context.Background(), "handlerQueryRowToMap error: %v", err)
			continue
		}
		maps = append(maps, toMap)
	}
	return maps
}
