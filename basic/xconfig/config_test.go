package xconfig

import (
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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
func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range m {
		strKey := fmt.Sprintf("%v", key)
		switch v := value.(type) {
		case map[interface{}]interface{}:
			result[strKey] = convertMap(v) // 递归处理嵌套 map
		case []interface{}:
			result[strKey] = convertSlice(v) // 处理 slice 中的嵌套结构
		default:
			result[strKey] = value
		}
	}
	return result
}

func convertSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, item := range s {
		switch v := item.(type) {
		case map[interface{}]interface{}:
			result[i] = convertMap(v) // 处理 slice 中的 map
		case []interface{}:
			result[i] = convertSlice(v) // 处理 slice 中的嵌套 slice
		default:
			result[i] = v
		}
	}
	return result
}

func TestYamlToJson(t *testing.T) {
	config := map[string]interface{}{}
	yamlContent := `
`
	err := yaml.Unmarshal([]byte(yamlContent), &config)
	if err != nil {
		fmt.Printf("unmarshal conf failed, err:%s \n", err)
		return
	}

	m2 := make(map[string]interface{})
	for key, val := range config {
		switch v := val.(type) {
		case map[interface{}]interface{}:
			m2[key] = convertMap(v) // 处理顶层 map
		case []interface{}:
			m2[key] = convertSlice(v) // 处理顶层 slice
		default:
			m2[key] = val
		}
	}

	indent, err := json.MarshalIndent(&m2, "", "  ")
	if err != nil {
		fmt.Printf("marshal indent failed, err:%s \n", err)
		return
	}
	fmt.Printf("%s\n", string(indent))
}
