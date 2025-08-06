package main

import (
	"example.com/xserver/biz"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xconfig/str_scanner"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth/aes_auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/cors"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/limiter"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xredis"
)

func main() {
	c := xconfig.Init(str_scanner.Default(`{"server":{"port":9981,"debug":true},"db":{"dns":"host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai","logger":true},"redis":{"addr":"127.0.0.1:6379","pass":"jerry123!","db":0}}`), "")
	engine := xe.New(c)
	engine.
		MustProvide(xdb.InitWith).
		MustProvide(xredis.InitWith).
		MustProvide(func() xe.EngineServer {
			server := xhs.New(engine)
			server.Use(limiter.Limiter(server.Config.Rate))
			server.Use(cors.New(server.Config))
			server.Use(auth.NewMw(aes_auth.New()))
			biz.Register(server)
			return server
		}).
		Start()
}
