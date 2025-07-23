package test

import (
	"context"
	"github.com/77d88/go-kit/basic/xconfig/json_scanner"
	"github.com/77d88/go-kit/basic/xrandom"
	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/external/xwx/wxmini"
	"github.com/77d88/go-kit/external/xwx/wxopen"
	"github.com/77d88/go-kit/external/xwx/wxpay"
	"github.com/77d88/go-kit/plugins/xapi"
	"github.com/77d88/go-kit/plugins/xapi/apis/admin_pro"
	"github.com/77d88/go-kit/plugins/xapi/auth"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xjob"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/plugins/xredis"
	"testing"
)

type To struct {
	To   string `json:"to"`
	Name string `json:"name"`
}

func TestTest(t *testing.T) {
	xlog.WithDebugger()
	xdb.Init(&xdb.Config{
		Dns:    "host=120.26.171.200 port=21001 user=postgres password=huanxi! dbname=hongyun sslmode=disable TimeZone=Asia/Shanghai",
		Logger: true,
	})
	xdb.Ctx(xapi.NewTestContext()).Exec("select 1")
	run := func(c *xhs.Ctx) {
		c.TraceId = -99
		//t.Log("c.Fatalf(\"any2\", xerror.newCtx(\"123\")) 最后一个参数是xError 不答应错误日志 并在返回msg中附带错误信息")
		//c.Fatalf("any2", xerror.newCtx("123"))
		//
		//t.Log(fmt.Sprintf("c.Fatalf(xerror.newCtx(\"123\"), \"用户端信息\", \"后台日志信息 %%s\", \"ss\") 不打印日志 返回自定义信息 并打印日志"))
		//
		//c.Fatalf(xerror.newCtx("123"), "用户端信息", "后台日志信息 %s", "ss")
		//
		//t.Log("c.Fatalf(xerror.newCtx(\"123\")) xerror 不打印日志 并在返回msg中附带错误信息")
		//c.Fatalf(xerror.newCtx("123"))
		//
		//t.Log(" c.Fatalf(xerror.newCtx(\"123\"), xapi.FatalWithCode(100)) xerror 不打印日志 并在返回msg中附带错误信息 同时设置错误码 这种情况需要使用xapi.FatalWithMsg自定义返回参数")
		//c.Fatalf(xerror.newCtx("123"), xapi.FatalWithCode(100))
		//
		//t.Log("c.Fatalf(\"123\")    不打印错误日志 并在返回msg中附带错误信息")
		//c.Fatalf("123")
		//
		//t.Log("c.Fatalf(errors.newCtx(\"123\"))    打印错误日志 并在返回info中附带错误信息")
		//c.Fatalf(errors.newCtx("123"))
		//
		//t.Log("c.Fatalf(errors.newCtx(\"123\"), \"错误\", \"执行\")  执行答应错误信息  并在msg中附带自定义消息")
		//c.Fatalf(errors.newCtx("123"), "错误", "执行")

		//t.Log("c.Fatalf(\"any\", errors.newCtx(\"123\")) 最后一个参数是error 则打印错误日志 并在返回info中附带错误信息")
		//c.Fatalf("any", errors.newCtx("123"))
		//
		//t.Log("c.Fatalf(nil, \"错误\")    不执行")
		//c.Fatalf(nil, "错误")
		//
		//t.Log("c.Fatalf(false, \"错误\")   不执行")
		//c.Fatalf(false, "错误")
		//
		//t.Log("c.Fatalf(true, \"用户端信息\", \"执行\") 执行打印错误信息  并在msg中附带消息")
		//c.Fatalf(true, "用户端信息", "执行")

		t.Log("c.Fatalf(xdb.Result{\n\t\t\tError:        gorm.ErrRecordNotFound,\n\t\t\tRowsAffected: 0,\n\t\t}) 执行打印错误信息  并在info中附带消息")
		//c.Fatalf(&xdb.Result{
		//	Error:        gorm.ErrRecordNotFound,
		//	RowsAffected: 0,
		//})

		c.Fatalf(true, "错误", "hh %s", "ss", xhs.FatalWithCode(-2))
		c.Fatalf(true, "错误", "hh %s", xhs.FatalWithCode(-2), "ss")

	}

	run(xapi.NewTestContext())
}

func TestTran(t *testing.T) {

}

func Txx(c context.Context) error {
	xdb.Ctx(c).Exec(`update s_live_room set popularity= ? where cid = 1`, xrandom.RandInt(1, 200))
	return xdb.CtxTran(c, func(d *xdb.DataSource) error {
		return xdb.Ctx(c).Exec(`update s_live_room set key= ? where id = 1`, xstr.ToString(xrandom.RandInt(1, 200))).Error

	})
}

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
		Add(wxmini.InitWith).
		Add(wxopen.InitWith).
		Add(wxpay.InitWith).
		Add(admin_pro.RegisterApi).
		AddAfter(func() auth.ICache {
			return xredis.NewUserBlackCache("black:")
		})
	api.Start()

}
