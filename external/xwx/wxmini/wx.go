package wxmini

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

var (
	Mini *miniprogram.MiniProgram
	Cfg  *miniConfig.Config
)

type RedisConfig struct {
	Addr string `yaml:"addr" json:"addr"` // 地址 ip:端口
	Pass string `yaml:"pass" json:"pass"` // 密码
	Db   int    `yaml:"db" json:"db"`     // 数据库
}

func InitWith(scanner xtype.Scanner, names ...string) *miniprogram.MiniProgram {
	//var redisConfig RedisConfig
	var config miniConfig.Config
	scanner.ScanKey(xarray.FirstOrDefault(names, "wx.mini"), &config)
	return Init(&config)
}

func Init(config *miniConfig.Config) *miniprogram.MiniProgram {
	Cfg = config
	if Cfg.AppID == "" {
		xlog.Errorf(nil, "wx mini xconfig is empty")
		return nil
	}
	//xconfig.ScanKey("redis", &redisConfig)
	//var c string
	//if redisConfig.Addr != "" {
	//	Cfg.Cache = cache.NewRedis(context.Background(), &cache.RedisOpts{
	//		Host:     redisConfig.Addr,
	//		Password: redisConfig.Pass,
	//		Database: redisConfig.Db,
	//	})
	//	c = "xredis"
	//} else {
	// 暂时全用本地了 没有多服务情况
	Cfg.Cache = cache.NewMemory()
	c := "memory"
	//}

	wc := wechat.NewWechat()
	Mini = wc.GetMiniProgram(Cfg)
	xlog.Infof(nil, "init wx mini success %s to xcache %s", Cfg.AppID, c)
	return Mini
}
