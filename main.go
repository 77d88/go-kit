package main

import (
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xconfig/json_scanner"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xe"
)

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
	api := xe.New(c).MustProvide(xhs.New)
	//api.Server.POST("/test", func(c *ctx.Ctx) {
	//	m := make(map[string]interface{})
	//	result := xdb.Ctx(c).Raw("select * from s_user limit 1").Scan(&m)
	//	c.Fatalf(result.GetError())
	//	c.Send(m)
	//})
	api.Start()
}
