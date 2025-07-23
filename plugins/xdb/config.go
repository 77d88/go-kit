package xdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xlog"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"time"
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

func InitWith(e *xe.Engine) *DataSource {
	var config Config
	e.Cfg.ScanKey(dbStr, &config)
	config.DbLinkName = dbStr
	return Init(&config)
}

func Init(c *Config) *DataSource {
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
		c.MaxOpenConns = 100
	}
	if c.MaxOpenConns == 0 {
		c.ConnMaxLifetime = 300
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 1800
	}
	if c.DbLinkName == "" {
		c.DbLinkName = dbStr

	}

	opts := &gorm.Config{
		QueryFields: true,
		Logger: &gormLogger{
			logger:        c.Logger,
			SlowThreshold: c.SlowThreshold,
		},
		CreateBatchSize: 1000, // 批量插入数量
	}

	gormDb, err := gorm.Open(postgres.Open(c.Dns), opts)
	if err != nil {
		xlog.Fatalf(nil, "db init error -1 %s", err)
		return nil
	}

	// 获取通用数据库对象 sql.DataSource ，然后使用其提供的功能
	sqlDB, err := gormDb.DB()
	if err != nil {
		xlog.Fatalf(nil, "db init error -2 %s", err)
		return nil
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)

	re := regexp.MustCompile(`password=.+? `)
	maskedStr := re.ReplaceAllString(c.Dns, "password=******* ")
	xlog.Infof(nil, "init db conn success %s link -> %s", maskedStr, c.DbLinkName)
	dbs[c.DbLinkName] = gormDb
	return &DataSource{
		DB:      gormDb,
		Context: context.TODO(),
		Config:  c,
	}
}

type gormLogger struct {
	logger        bool
	SlowThreshold int
}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	xlog.Infof(ctx, s, i...)
}

func (g gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	xlog.Warnf(ctx, s, i...)
}

func (g gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	xlog.Errorf(ctx, s, i...)
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if !g.logger {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	// 把sql 的换行改为 \t
	sql = regexp.MustCompile(`\n|\r\n|\r`).ReplaceAllString(sql, "\t")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		xlog.DefaultLogger.WithOptions(zap.AddCallerSkip(4)).Debug(fmt.Sprintf("sql err: %s \n %v", sql, err), xlog.GetFields(ctx)...)
	} else {
		if elapsed > time.Duration(g.SlowThreshold)*time.Millisecond {
			xlog.DefaultLogger.WithOptions(zap.AddCallerSkip(2)).Debug(fmt.Sprintf("[%.2fms:%d] query slow %s", float64(elapsed.Nanoseconds())/1e6, rows, sql), xlog.GetFields(ctx)...)
		} else {
			xlog.DefaultLogger.WithOptions(zap.AddCallerSkip(2)).Debug(fmt.Sprintf("[%.2fms:%d]: %s", float64(elapsed.Nanoseconds())/1e6, rows, sql), xlog.GetFields(ctx)...)
		}
	}
}
