package xconfig

import (
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xsys"
	"github.com/spf13/viper"
)

type Config struct {
	viper         *viper.Viper
	loader        ConfigLoader      // 配置加载器
	cacheConfig   map[string]string // 缓存配置
	dataIds       []string          // 配置文件
	ListenDataIds []string          // 监听的配置文件
	listenStop    chan struct{}
	group         string
}

var XConfig *Config

// ConfigLoader 配置加载器
type ConfigLoader interface {
	Load(group, dataId string) (string, error)
	Type() string
}

func Init(loader ConfigLoader, group string, dataIds ...string) *Config {

	baseConfig := xsys.OsEnvGet("V_DEFAULT_CONFIG_KEY", "base")
	if baseConfig != "ignore" {
		var bdataIds = make([]string, 1)
		bdataIds[0] = baseConfig
		dataIds = append(bdataIds, dataIds...)
	}
	dataIds = xarray.Union(dataIds) // 去重

	if len(dataIds) == 0 {
		panic("dataIds is empty")
	}

	cfg := Config{
		group:       group,
		loader:      loader,
		cacheConfig: make(map[string]string),
		dataIds:     dataIds,
		listenStop:  make(chan struct{}),
	}
	cfg.readToViper()
	cfg.startListen() // 开启监听
	InfoLog("config init success from to %s : [%s]==> successful:%v error:%v all:%v", loader.Type(), cfg.group, cfg.ListenDataIds, xarray.Difference(cfg.dataIds, cfg.ListenDataIds), cfg.dataIds)
	if len(cfg.ListenDataIds) > 0 {
	} else {
		WarnLog("no config to listen input %s->%v", cfg.group, dataIds)
	}
	XConfig = &cfg
	return &cfg
}

func Scan(config any) {
	if err := XConfig.viper.Unmarshal(config); err != nil {
		fmt.Printf("unmarshal conf failed, err:%s \n", err)
	}
}

func ScanKey(key string, config any) {
	if err := XConfig.viper.UnmarshalKey(key, config); err != nil {
		fmt.Printf("unmarshal conf key %s failed, err:%s \n", key, err)
	}
}

func GetString(key string) string {
	return XConfig.viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return XConfig.viper.GetStringSlice(key)
}

func ScanKeyExecute(key string, config any, f func()) {
	if err := XConfig.viper.UnmarshalKey(key, config); err != nil {
		fmt.Printf("unmarshal conf key %s failed, err:%s \n", key, err)
	} else {
		f()
	}
}


func (c *Config) Dispose() error {
	c.listenStop <- struct{}{}
	return nil
}


// ShutDownListen 关闭配置监听
func (c *Config) ShutDownListen() {
	c.listenStop <- struct{}{}
}

func (c *Config) Scan(config any) {
	if err := c.viper.Unmarshal(config); err != nil {
		fmt.Printf("unmarshal conf failed, err:%s \n", err)
	}
}

func (c *Config) ScanKey(key string, config any) {
	if err := c.viper.UnmarshalKey(key, config); err != nil {
		fmt.Printf("unmarshal conf key %s failed, err:%s \n", key, err)
	}
}

func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *Config) ScanKeyExecute(key string, config any, f func()) {
	if err := c.viper.UnmarshalKey(key, config); err != nil {
		fmt.Printf("unmarshal conf key %s failed, err:%s \n", key, err)
	} else {
		f()
	}
}
