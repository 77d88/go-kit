package main

import (
	"example.com/xserver/biz"
	"github.com/77d88/go-kit/basic/xconfig/str_scanner"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth/redis_auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/cors"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/limiter"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xdatabase/xredis"
	"github.com/77d88/go-kit/plugins/xdatabase/xredis/redismq"
	"github.com/77d88/go-kit/plugins/xtask/xcron"
	"github.com/77d88/go-kit/plugins/xtask/xjob"
)

func init() {
	str_scanner.Default(`{"server":{"port":9981,"debug":true},"db":{"dns":"host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai","logger":true},"redis":{"addr":"127.0.0.1:6666","pass":"test"}}`)
}

func main() {
	// 与初始化
	x.FastInit(func(*xdb.DB, *xredis.Client, *xjob.Manager, *xcron.Manager,*redismq.Client) {})
	x.Use(func() x.EngineServer {
		hs := xhs.New().Use(limiter.New()).Use(cors.New()).Use(auth.NewMw(redis_auth.New()))
		biz.Register(hs)
		return hs
	})

	x.Start()

}
