package xapi

import (
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xconfig/json_scanner"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth/aes_auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/cors"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/plugins/xredis"
	"testing"
)

func a(c *xhs.Ctx, db *xdb.DataSource, redis *xredis.RedisClient) (interface{}, error) {
	m := make(map[string]interface{})
	xdb.Ctx(c).Raw("select * from s_user limit 1").Scan(&m)
	db.WithContext(c).Raw("select * from s_user limit 1").Scan(&m)
	xlog.Errorf(c, "%v", m)
	return m, nil
}

func TestName(t *testing.T) {

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
			server.Use(cors.New(server.Config))
			server.Use(auth.NewMw(aes_auth.New()))
			server.GET("/test", server.WrapWithDI(a))
			return server
		}).
		Start()

}
