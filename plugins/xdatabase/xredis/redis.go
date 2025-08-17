package xredis

import (
	"context"
	"errors"
	"sync"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil
const redisStr string = "redis"

var (
	dbs    = make(map[string]*redis.Client)
	dbLock = sync.RWMutex{} // 添加读写锁保护 dbs map
)

// Config redis 配置
type Config struct {
	Addr       string `yaml:"addr" json:"addr"` // 地址 ip:端口
	Pass       string `yaml:"pass" json:"pass"` // 密码
	Db         int    `yaml:"db" json:"db"`     // 数据库
	DbLinkName string `yaml:"dbLinkName" json:"dbLinkName"`
}

// client redis命令
type client struct {
	*redis.Client
	DbLinkName string
	Config     *Config
}

func (c *client) Dispose() error {
	err := c.Close()
	if err != nil {
		xlog.Errorf(nil, "redis dispose errror link->%s ", c.DbLinkName)
	}
	xlog.Infof(nil, "redis dispose success<%s:%d> link->%s ", c.Config.Addr, c.Config.Db, c.DbLinkName)
	return err
}

func NewX() *redis.Client {
	c, err := x.Config[Config](redisStr)
	if err != nil {
		xlog.Errorf(nil, "redis init Fatal %+v", err)
		return nil
	}
	c.DbLinkName = redisStr
	return New(c)
}

func New(config *Config) *redis.Client {

	if config == nil {
		xlog.Errorf(nil, "redis init Fatal %+v", config)
		return nil
	}
	if config.Addr == "" {
		xlog.Errorf(nil, "redis init Fatal addr %+v", config)
		return nil
	}
	// 使用读锁检查是否已经存在相同名称的连接
	dbLock.RLock()
	existingClient, exists := dbs[config.DbLinkName]
	dbLock.RUnlock()
	// 如果已经存在相同名称的连接，直接返回已存在的实例
	if exists {
		xlog.Debugf(nil, "redis connection with name %s already exists, returning existing instance", config.DbLinkName)
		return existingClient
	}

	// 使用写锁确保只有一个 goroutine 能够创建连接
	dbLock.Lock()
	defer dbLock.Unlock()

	Client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pass, // no password set
		DB:       config.Db,   // use default Database
	})
	cli := client{
		Client:     Client,
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
	dbs[config.DbLinkName] = Client
	x.Use(&cli, "__xredis."+config.DbLinkName)
	return Client

}

// Get 获取数据库链接
func Get(name ...string) (*redis.Client, error) {
	dbLock.RLock()
	defer dbLock.RUnlock()

	database, ok := dbs[xarray.FirstOrDefault(name, redisStr)]
	if !ok {
		xlog.Errorf(nil, "数据库[%s]链接不存在", name)
		return nil, xerror.New("数据库链接不存在")
	}
	return database, nil
}

// GetAll 获取所有数据库连接
func GetAll() map[string]*redis.Client {
	dbLock.RLock()
	defer dbLock.RUnlock()

	// 返回副本以避免外部修改
	result := make(map[string]*redis.Client, len(dbs))
	for k, v := range dbs {
		result[k] = v
	}
	return result
}
