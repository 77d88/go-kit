package xpg

import (
	"context"
	"regexp"

	"github.com/77d88/go-kit/plugins/xlog"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB
type DB struct {
	pool   *pgxpool.Pool
	config *Config
}

func C(c context.Context, name ...string) *Inst {
	db, err := GetDB(name...)
	if err != nil {
		return nil
	}
	return &Inst{
		pool:   db.pool,
		cond:   sq.Select().PlaceholderFormat(sq.Dollar),
		ctx:    c,
		config: db.config,
	}
}
func (db *DB) clone() *DB {
	return &DB{
		pool:   db.pool,
		config: db.config,
	}
}

func (db *DB) Dispose() error {
	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(db.config.Dns, "password=******* ")
	xlog.Warnf(nil, "close db conn %s<%s>", db.config.DbLinkName, maskedStr)
	db.pool.Close()
	return nil
}
