package xpg

import (
	"errors"
	"fmt"
	"time"

	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

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
	txInst := i.Copy()
	txInst.cond = sq.SelectBuilder{}
	txInst.tx = begin

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
