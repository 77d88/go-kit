package xconfig

import (
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/spf13/viper"
	"strings"
	"testing"
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

func Test(t *testing.T) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("v")                                // 设置读取环境变量前缀，会自动转为大写 v_APP_CLIENT_ID=1 v_db.dns =
	viper.AllowEmptyEnv(true)                              // 将空环境变量视为已设置
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 环境变量中配置的.替换为_书写

	fmt.Printf("%s\n", viper.GetString("NACOS_ADDR"))
	fmt.Printf("%s\n", viper.GetString("NACOS_NAMESPACE"))
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

