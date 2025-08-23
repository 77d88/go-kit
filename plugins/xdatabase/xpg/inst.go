package xpg

import (
	"context"
	"sync"

	"github.com/77d88/go-kit/basic/xreflect"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Inst struct {
	pool             *pgxpool.Pool
	config           *Config
	cond             sq.SelectBuilder
	ctx              context.Context
	tx               pgx.Tx
	savepointCounter int
	spMu             sync.Mutex
	txMarkedBad      bool
	debug            bool
	selectFields     []string
	tableName        string
	model            Model
}

func (i *Inst) Debug() *Inst {
	i.debug = true
	return i
}

// WithContext 设置上下文
func (i *Inst) WithContext(ctx context.Context) *Inst {
	inst := i.Copy()
	inst.ctx = ctx
	return inst
}

// Copy 复制 所有都一样包括事务都一样
func (i *Inst) Copy() *Inst {
	return &Inst{
		pool:             i.pool,
		config:           i.config,
		cond:             i.cond,
		ctx:              i.ctx,
		tx:               i.tx,
		savepointCounter: i.savepointCounter,
		txMarkedBad:      i.txMarkedBad,
		debug:            i.debug,
		selectFields:     i.selectFields,
		tableName:        i.tableName,
		model:            i.model,
	}
}

func (i *Inst) Model(model Model) *Inst {
	inst := i.Copy()
	inst.model = model
	inst.tableName = model.TableName()
	value, _ := xreflect.GetInst(model)
	fields := extractDBFields(value.Type())
	inst.selectFields = fields
	return inst
}

func (i *Inst) Table(table string) *Inst {
	inst := i.Copy()
	inst.tableName = table
	return inst
}

func (i *Inst) Select(fields ...string) *Inst {
	inst := i.Copy()
	inst.selectFields = fields
	return inst
}

func (i *Inst) ToSql() (string, []interface{}, error) {
	return i.cond.From(i.tableName).Columns(i.selectFields...).ToSql()
}
