package xredis

import (
	"context"
	"errors"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil
const redisStr string = "redis"


var (
	dbs = make(map[string]*Client)
)

// Config redis 配置
type Config struct {
	Addr       string `yaml:"addr" json:"addr"` // 地址 ip:端口
	Pass       string `yaml:"pass" json:"pass"` // 密码
	Db         int    `yaml:"db" json:"db"`     // 数据库
	DbLinkName string `yaml:"dbLinkName" json:"dbLinkName"`
}

// Client redis命令
type Client struct {
	*redis.Client
	DbLinkName string
	Config     *Config
}

func (c *Client) Dispose() error {
	err := c.Close()
	if err != nil {
		xlog.Errorf(nil, "redis dispose errror link->%s ", c.DbLinkName)
	}
	xlog.Infof(nil, "redis dispose success<%s:%d> link->%s ", c.Config.Addr, c.Config.Db, c.DbLinkName)
	return err
}

func InitWith(e *xe.Engine) *Client {
	var config Config
	e.Cfg.ScanKey(redisStr, &config)
	config.DbLinkName = redisStr
	return Init(&config)
}

func Init(config *Config) *Client {

	if config == nil {
		xlog.Errorf(nil, "redis init Fatal %+v", config)
		return nil
	}
	if config.Addr == "" {
		xlog.Errorf(nil, "redis init Fatal addr %+v", config)
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pass, // no password set
		DB:       config.Db,   // use default Database
	})
	cli := Client{
		Client:     client,
		DbLinkName: config.DbLinkName,
		Config:     config,
	}
	// 测试连接
	cmd := cli.Get(context.Background(), "##test")
	if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
		xlog.Errorf(nil, "redis init Fatal %s db %d %+v link ->%s", config.Addr, config.Db, cmd.Err(), config.DbLinkName)
		return nil
	} else {
		xlog.Infof(nil, "redis init success %s db %d -> %s", config.Addr, config.Db, config.DbLinkName)
	}
	dbs[config.DbLinkName] = &cli
	return &cli

}

// Get 获取数据库链接
func Get(name ...string) (*Client, error) {
	database, ok := dbs[xarray.FirstOrDefault(name, redisStr)]
	if !ok {
		xlog.Errorf(nil, "数据库[%s]链接不存在", name)
		return nil, xerror.New("数据库链接不存在")
	}
	return database, nil
}
