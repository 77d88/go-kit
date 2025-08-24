package xpg

import (
	"context"
	"regexp"
	"time"

	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config 数据库配置
type Config struct {
	Dns                       string `json:"dns"`                       // 示例 host=127.0.0.1 port=5432 user=postgres password=yzz123! dbname=yzz sslmode=disable TimeZone=Asia/Shanghai
	Logger                    bool   `json:"logger"`                    // 是否打印sql日志
	MaxIdleConns              int    `json:"maxIdleConns"`              // MaxIdleConns 用于设置连接池中空闲连接的最大数量。
	MaxOpenConns              int    `json:"maxOpenConns"`              // MaxOpenConns 设置打开数据库连接的最大数量。
	ConnMaxLifetime           int    `json:"connMaxLifetime"`           // ConnMaxLifetime 设置了连接可复用的最大时间。 单位 s
	SlowThreshold             int    `json:"slowThreshold"`             // 慢查询阈值 默认500ms
	IgnoreRecordNotFoundError *bool  `yaml:"ignoreRecordNotFoundError"` // 忽略记录未找到的错误
	DbLinkName                string `yaml:"dbLinkName"`                // 数据库链接名称
}

func NewX() *DB {
	c, err := x.Config[Config](DefaultDbLinkStr)
	if err != nil {
		xlog.Fatalf(nil, "init db error %s", err)
		return nil
	}
	c.DbLinkName = DefaultDbLinkStr
	return New(c)
}

func New(c *Config) *DB {
	if c == nil {
		xlog.Fatalf(nil, "init db error config is nil ")
		return nil
	}

	if len(c.Dns) == 0 {
		xlog.Fatalf(nil, "init %s db error config.dns is nil ", c.DbLinkName)
		return nil
	}

	if c.SlowThreshold == 0 {
		c.SlowThreshold = 500
	}
	if c.IgnoreRecordNotFoundError == nil {
		c.IgnoreRecordNotFoundError = xcore.V2p(true)
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}
	if c.MaxOpenConns == 0 {
		c.ConnMaxLifetime = 300
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 1800
	}
	if c.DbLinkName == "" {
		c.DbLinkName = DefaultDbLinkStr

	}

	ctx := context.Background()
	config, err := pgxpool.ParseConfig(c.Dns)
	if err != nil {
		xlog.Fatalf(nil, "init db error %s", err)
		return nil
	}
	// 最大连接数
	config.MaxConns = int32(c.MaxOpenConns)
	config.MinConns = 2
	// 空闲连接的最大存活时间，超过此时间的空闲连接会在健康检查时被关闭。
	config.MaxConnIdleTime = time.Duration(c.MaxIdleConns) * time.Second
	// 连接的最大存活时间，超过此时间的连接会被自动关闭。
	config.MaxConnLifetime = time.Duration(c.ConnMaxLifetime) * time.Second
	// 在 MaxConnLifetime 基础上增加的随机时间，防止所有连接同时过期。
	config.MaxConnLifetimeJitter = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		xlog.Fatalf(nil, "db init error -3 %s", err)
		return nil
	}

	go func() {
		err = pool.Ping(ctx)
		if err != nil {
			xlog.Fatalf(nil, "db init error -1 %s", err)
		}
	}()

	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(c.Dns, "password=******* ")
	xlog.Infof(nil, "init db conn success %s link -> %s", maskedStr, c.DbLinkName)
	db := &DB{
		pool:   pool,
		config: c,
	}
	dbs[c.DbLinkName] = db
	if DefaultDB == nil {
		DefaultDB = db
	}
	return db
}
