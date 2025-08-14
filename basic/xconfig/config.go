package xconfig

import (
	"fmt"
	"sync"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/spf13/viper"
)

type Config struct {
	viper         *viper.Viper
	loader        ConfigLoader      // 配置加载器
	cacheConfig   map[string]string // 缓存配置
	dataIds       []string          // 配置文件
	ListenDataIds []string          // 监听的配置文件
	listenStop    chan struct{}
}

var XConfig *Config
var once sync.Once

// ConfigLoader 配置加载器
type ConfigLoader interface {
	Load(dataId string) (string, error)
	Type() string
}

func Init(loader ConfigLoader, dataIds ...string) *Config {
	once.Do(func() {
		dataIds = xarray.Union(dataIds) // 去重
		cfg := Config{
			loader:      loader,
			cacheConfig: make(map[string]string),
			dataIds:     dataIds,
			listenStop:  make(chan struct{}),
		}
		cfg.readToViper()
		cfg.startListen() // 开启监听
		InfoLog("config init success from to %s ==> successful:%v error:%v all:%v", loader.Type(), cfg.ListenDataIds, xarray.Difference(cfg.dataIds, cfg.ListenDataIds), cfg.dataIds)
		if len(cfg.ListenDataIds) > 0 {
		} else {
			WarnLog("no config to listen input %v", dataIds)
		}
		XConfig = &cfg
	})
	return XConfig
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

func (c *Config) Scan(config any) error {
	if err := c.viper.Unmarshal(config); err != nil {
		ErrorLog("unmarshal conf failed, err:%s \n", err)
		return xerror.Newf("unmarshal conf failed, err:%s \n", err)
	}
	return nil
}

func (c *Config) ScanKey(key string, config any) error {
	if err := c.viper.UnmarshalKey(key, config); err != nil {
		ErrorLog("unmarshal conf key %s failed, err:%s \n", key, err)
		return xerror.Newf("unmarshal conf key %s failed, err:%s \n", key, err)
	}
	return nil
}

func (c *Config) GetString(key string, defaultValue ...string) string {
	str := c.viper.GetString(key)
	if str == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return str
}

func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *Config) GetInt(key string, defaultValue ...int) int {
	i := c.viper.GetInt(key)
	if i == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
	}
	return i
}

func (c *Config) GetIntSlice(key string) []int {
	return c.viper.GetIntSlice(key)
}

func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *Config) ScanKeyExecute(key string, config any, f func()) {
	if err := c.viper.UnmarshalKey(key, config); err != nil {
		fmt.Printf("unmarshal conf key %s failed, err:%s \n", key, err)
	} else {
		f()
	}
}
