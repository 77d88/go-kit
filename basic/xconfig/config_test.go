package xconfig

import (
	"fmt"
	"testing"

	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/spf13/viper"
)

type (
	ServerConfig struct {
		Port int `yaml:"port"`
		Rate int `yaml:"rate"`
	}

	DbConfig struct {
		Dns                       string `json:"dns"`                       // 示例 host=127.0.0.1 port=5432 user=postgres password=yzz123! dbname=yzz sslmode=disable TimeZone=Asia/Shanghai
		Logger                    bool   `json:"logger"`                    // 是否打印sql日志
		MaxIdleConns              int    `json:"MaxIdleConns"`              // MaxIdleConns 用于设置连接池中空闲连接的最大数量。
		MaxOpenConns              int    `json:"MaxOpenConns"`              // MaxOpenConns 设置打开数据库连接的最大数量。
		ConnMaxLifetime           int    `json:"ConnMaxLifetime"`           // ConnMaxLifetime 设置了连接可复用的最大时间。 单位 s
		SlowThreshold             int    `json:"SlowThreshold"`             // 慢查询阈值 默认500ms
		IgnoreRecordNotFoundError *bool  `json:"IgnoreRecordNotFoundError"` // 忽略记录未找到的错误
	}
)

type StringLoader struct {
	data string
}

func (c *StringLoader) Load(group, dataId string) (string, error) {
	return c.data, nil
}

func (c *StringLoader) Type() string {
	return "static json string"
}

func Test(t *testing.T) {
	config := Init(&StringLoader{data: `{"server":{"port":9981,"debug":false},"db":{"dns":"host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai","logger":true}}`},"")


	fmt.Printf("%d\n", config.GetInt("server.port",8080)) // 环境变量设置 x_server.port=80 输出80 不设置输出 9981. 配置值
	fmt.Printf("%s\n", config.GetString("db.dns")) // 环境变量设置 x_db.dns=tesxx 则输出tesxx 不设置输出 host.... 配置值
	fmt.Printf("%s\n", config.GetString("test")) // 环境变量设置为x_test=x_test_value
	fmt.Printf("%s\n", viper.GetString("timeZone"))
	fmt.Printf("appid %s\n", viper.GetString("wx.mini.appid"))
}

func TestConfig(t *testing.T) {
	xlog.WithDebugger()
	fmt.Printf("%s\n", viper.GetString("server.port"))
	fmt.Printf("%s\n", viper.GetString("server.name"))
	fmt.Printf("%s\n", viper.GetString("server.rate"))
	fmt.Printf("%s\n", viper.GetString("wx.features.getPhoneNumber"))

}

