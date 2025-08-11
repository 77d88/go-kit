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
	"github.com/77d88/go-kit/plugins/xtask/xcron"
	"github.com/77d88/go-kit/plugins/xtask/xjob"
)

func main() {
	str_scanner.Default(`{"server":{"port":9981,"debug":true},"db":{"dns":"host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai","logger":true},"redis":{"addr":"127.0.0.1:6379","pass":"jerry123!","db":0}}`)
	x.Use(xdb.InitWith)
	x.Use(xredis.InitWith)
	x.Use(xjob.Init)
	x.Use(xcron.Init)
	x.Use(func() x.EngineServer {
		server := xhs.New()
		server.Use(limiter.Limiter(server.Config.Rate))
		server.Use(cors.New(server.Config))
		server.Use(auth.NewMw(redis_auth.New()))
		biz.Register(server)
		return server
	})
	x.Start()

}
