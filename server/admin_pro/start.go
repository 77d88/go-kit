package main

import (
	"github.com/77d88/go-kit/basic/xconfig/str_scanner"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/auth/aes_auth"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/cors"
	"github.com/77d88/go-kit/plugins/xapi/server/mw/limiter"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xjob"
	"github.com/77d88/go-kit/plugins/xredis"
	"github.com/77d88/go-kit/server/admin_pro/proapi"
	"github.com/77d88/go-kit/server/xaliyun/aliyunaddress"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossgin"
	"github.com/77d88/go-kit/server/xaliyun/aliyunoss/aliyunossginsts"
)

func main() {
	sc := str_scanner.Default(`{
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
	}`)
	xe.New(sc).
		Use(xdb.InitWith).
		Use(xredis.InitWith).
		Use(xjob.Init).
		Use(aliyunoss.InitWith, true).
		UseServer(func(e *xe.Engine) (xe.EngineServer,error) {
			server := xhs.New(e)
			server.Use(limiter.Limiter(server.Config.Rate))
			server.Use(cors.New(server.Config))
			server.Use(auth.NewMw(aes_auth.New()))
			proapi.Register(server)
			aliyunossginsts.DefaultRegister("/aliyun/sts", server)
			aliyunaddress.DefaultRegister("/aliyun/address", server)
			aliyunossgin.DefaultRegister("/aliyun/oss", server)
			return server,nil
		}).Start()
}
