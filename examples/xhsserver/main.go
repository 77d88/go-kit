package main

import (
	"context"
	"example.com/xserver/biz"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xconfig/json_scanner"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth/aes_auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/cors"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/limiter"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xredis"
)

func c2(c context.Context, q string) {
	xdb.BeginWithCtx(c).Exec("update s_user set sys_nickname = ? where id = 600075249287237", q)
}
func a(c *xhs.Ctx, db *xdb.DB, redis *xredis.Client) (interface{}, error) {
	query := c.DefaultQuery("name", "test")

	c2(c, query)
	db.WithCtx(c).Exec("update s_user set note = ? where id = 600075249287237", query)
	if query == "test" {
		return nil, xerror.New("error no query")
	}
	return nil, nil
}

func b(c *xhs.Ctx, db *xdb.DB) (interface{}, error) {
	m := make(map[string]interface{})
	scan := db.WithCtx(c).Table("s_user").WithId(600075249287237).Scan(&m)
	return m, scan.Error
}

func main() {

	c := xconfig.Init(json_scanner.Default(`{
  "server": {
    "port": 9981,
    "debug": true,
	"logLevel": -1
  },
  "db": {
    "dns": "host=127.0.0.1 port=5432 user=postgres password=jerry123! dbname=zyv2 sslmode=disable TimeZone=Asia/Shanghai",
    "logger": true
  },
  "redis": {
    "addr": "127.0.0.1:6379",
    "pass": "jerry123!",
    "db": 0
  }
}`), "")
	engine := xe.New(c)
	engine.
		MustProvide(xdb.InitWith).
		MustProvide(xredis.InitWith).
		MustProvide(func() xe.EngineServer {
			server := xhs.New(engine)
			engine.MustProvide(func() *xhs.HttpServer {
				return server
			})
			server.Use(limiter.Limiter(server.Config.Rate))
			server.Use(cors.New(server.Config))
			server.Use(auth.NewMw(aes_auth.New()))
			biz.Register(server)
			return server
		}).
		Start()
}
