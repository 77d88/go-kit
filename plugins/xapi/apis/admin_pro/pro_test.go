package admin_pro

import (
	"github.com/77d88/go-kit/basic/xconfig/json_scanner"
	"github.com/77d88/go-kit/plugins/xapi/auth"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xjob"
	"github.com/77d88/go-kit/plugins/xredis"
	"testing"
)

func TestRun(t *testing.T) {
	sc := json_scanner.Default(`{
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
	api := xe.New(sc).
		Add(xdb.InitWith).
		Add(xredis.InitWith).
		Add(xjob.Init).
		Add(RegisterApi).
		AddAfter(func() auth.ICache {
			return xredis.NewUserBlackCache("black:")
		})
	api.Start()

}
