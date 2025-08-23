package main

import (
	"github.com/77d88/go-kit/basic/xconfig/str_scanner"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth/redis_auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/cors"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/limiter"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xpg"
	"github.com/77d88/go-kit/plugins/xdatabase/xredis"
	_ "github.com/77d88/go-kit/plugins/xdatabase/xredis"
	"github.com/77d88/go-kit/plugins/xtask/xjob"
	"github.com/77d88/go-kit/server/admin_pro/proapi"
)

func init() {
	x.Must(str_scanner.Default(`{
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
	   "addr": "127.0.0.1:6666",
	   "pass": "test",
	   "db": 0
	 }
	}`))
}

func main() {
	//x.Must(xdb.NewX)
	x.Must(xredis.NewX)
	x.Must(xjob.NewX)
	x.Must(redis_auth.New)
	x.Must(xpg.NewX)
	x.Use(func() x.EngineServer {
		server := xhs.New()
		server.Use(cors.New())
		server.Use(limiter.New())
		server.Use(auth.NewX().TokenInfo())
		proapi.Register(server)
		//aliyunossginsts.DefaultRegister("/aliyun/sts", server)
		//aliyunaddress.DefaultRegister("/aliyun/address", server)
		//aliyunossgin.DefaultRegister("/aliyun/oss", server)
		return server
	})
	x.Start()
}
