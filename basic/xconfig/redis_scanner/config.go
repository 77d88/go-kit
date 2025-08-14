package redis_scanner

import (
	"context"
	"errors"
	"fmt"

	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xdbutil"
	"github.com/77d88/go-kit/basic/xsys"
	"github.com/redis/go-redis/v9"
)

type RedisConfigLoader struct {
	configClient *redis.Client
	redisPrefix  string
	cnn          *xdbutil.ConnectionInfo
}

func Default(dataIds ...string) *xconfig.Config {
	config := xconfig.Init(NewEnv(), dataIds...)
	return config
}

func New(conn *xdbutil.ConnectionInfo, prefix string) *RedisConfigLoader {
	config := &RedisConfigLoader{
		redisPrefix: prefix,
		cnn:         conn,
	}
	config.configClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conn.Host, conn.Port),
		Password: conn.Password,        // no password set
		DB:       conn.ToIntDatabase(), // use default Database
		Username: conn.Username,
	})

	// 测试连接
	cmd := config.configClient.Get(context.Background(), "##test")
	if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
		xconfig.ErrorLog("config init Fatal %s db %d %+v", conn.Host, conn.ToIntDatabase(), cmd.Err())
		panic("redis init Fatal")
	}
	return config

}

// NewEnv 通过环境变量创建一个RedisConfig
func NewEnv() *RedisConfigLoader {
	redisPrefix := xsys.OsEnvGet("V_CONFIG_REDIS_PREFIX", "config:")
	redisAddr := xsys.OsEnvGet("V_CONFIG_REDIS_ADDR", "redis://default:jerry123!@127.0.0.1:6379/0")
	cnn, err := xdbutil.ParseConnection(redisAddr)
	if err != nil {
		xconfig.ErrorLog("config init Fatal %s %+v", redisAddr, err)
		panic(err)
	}
	return New(cnn, redisPrefix)

}
func (c *RedisConfigLoader) Load(dataId string) (string, error) {
	sprintf := fmt.Sprintf("%s:%s", c.redisPrefix, dataId)
	str, err := c.configClient.Get(context.TODO(), sprintf).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		xconfig.ErrorLog("config get %s  error : %s", dataId, err)
		return "", err
	}
	return str, nil
}
func (c *RedisConfigLoader) Type() string {
	return fmt.Sprintf("redis(%s:%d/%d)", c.cnn.Host, c.cnn.Port, c.cnn.ToIntDatabase())
}
