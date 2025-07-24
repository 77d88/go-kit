package wxopen

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xtype"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

type RedisConfig struct {
	Addr string `yaml:"addr" json:"addr"` // 地址 ip:端口
	Pass string `yaml:"pass" json:"pass"` // 密码
	Db   int    `yaml:"db" json:"db"`     // 数据库
}

var (
	Cfg      *offConfig.Config
	Official *officialaccount.OfficialAccount
)

func InitWith(scanner xtype.Scanner, names ...string) *officialaccount.OfficialAccount {
	var config offConfig.Config
	scanner.ScanKey(xarray.FirstOrDefault(names, "wx.open"), &config)
	return Init(&config)
}
func Init(cfg *offConfig.Config) *officialaccount.OfficialAccount {
	//var redisConfig RedisConfig
	Cfg = cfg
	if Cfg.AppID == "" {
		xlog.Errorf(nil, "wx official xconfig is empty")
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
	Cfg.Cache = cache.NewMemory()
	c := "memory"
	//}

	wc := wechat.NewWechat()
	Official = wc.GetOfficialAccount(Cfg)
	//
	xlog.Infof(nil, "init wx official success %s to xcache %s", Cfg.AppID, c)
	return Official
}
