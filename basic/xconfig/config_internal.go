package xconfig

import (
	"strings"
	"time"

	"github.com/77d88/go-kit/basic/xsys"
	"github.com/spf13/viper"
)

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
					config, err := c.loader.Load(id)
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

func (c *Config) readToViper() {
	c.viper = viper.New()         // 重置viper
	c.viper.SetConfigType("json") // 统一使用json格式
	// 添加环境变量支持
	c.viper.AutomaticEnv() // 自动绑定环境变量
	// 可选：设置环境变量前缀，例如 "APP_"
	c.viper.SetEnvPrefix("x")
	// 可选：替换环境变量中的 "." 为 "_"（Viper 默认将 "." 替换为 "_"）
	// c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	successKeys := make([]string, 0)
	for _, id := range c.dataIds {
		config, err := c.loader.Load(id)

		if err != nil {
			continue
		}
		if config == "" {
			continue
		}
		err = c.viper.MergeConfig(strings.NewReader(config)) // 读取配置文件到viper
		if err != nil {
			ErrorLog("viper read config 【%s】 error:%s", id, err)
			continue
		}
		c.cacheConfig[id] = config
		successKeys = append(successKeys, id)
	}
	c.ListenDataIds = successKeys
}
