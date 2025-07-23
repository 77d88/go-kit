package xconfig

import (
	"fmt"
	"github.com/77d88/go-kit/basic/xsys"
	"github.com/spf13/viper"
	"strings"
	"time"
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

// startListen 开启配置监听 每分钟 监听配置文件变化
func (c *Config) startListen() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ErrorLog("listen config panic:%v", r)
			}
		}()
		if len(c.dataIds) == 0 {
			return
		}
		// 每分钟执行一次
		//InfoLog("start listen config change %v", c.ListenDataIds)
		ticker := time.NewTicker(time.Minute * 1)
		for {
			select {
			case <-ticker.C:
				for _, id := range c.dataIds { // 监听所有数据源
					config, err := c.loader.Load(c.group, id)
					if err != nil {
						if err.Error() != "config data not exist" && id != "base" {
							ErrorLog("sync xconfig error:%s", err)
						}
						continue
					}
					// 不一致才更新
					if config != c.cacheConfig[id] {
						WarnLog("xconfig change %s", id)
						//readToViper(xconfig) 现在暴力一点 直接重启程序 反正也快
						xsys.Restart()
					}
				}

			case <-c.listenStop:
				WarnLog("stop listen config change")
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Config) Init() error {
	return nil
}

func (c *Config) Dispose() error {
	c.listenStop <- struct{}{}
	return nil
}

func (c *Config) readToViper() {
	c.viper = viper.New()         // 重置viper
	c.viper.SetConfigType("json") // 统一使用json格式
	successKeys := make([]string, 0)
	for _, id := range c.dataIds {
		config, err := c.loader.Load(c.group, id)

		if err != nil {
			continue
		}
		if config == "" {
			continue
		}
		err = c.viper.MergeConfig(strings.NewReader(config)) // 读取配置文件到viper
		if err != nil {
			ErrorLog("viper read config 【%s:%s】 error:%s", c.group, id, err)
			continue
		}
		c.cacheConfig[id] = config
		successKeys = append(successKeys, id)
	}
	c.ListenDataIds = successKeys
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
