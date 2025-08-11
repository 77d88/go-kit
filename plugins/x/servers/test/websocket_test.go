package main

//import (
//	"fmt"
//	"github.com/77d88/go-kit/plugins/xapi"
//	"github.com/77d88/go-kit/plugins/xapi/ctx"
//	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
//	"github.com/77d88/go-kit/plugins/xlog"
//	"github.com/77d88/go-kit/plugins/xqueue"
//	"github.com/77d88/go-kit/plugins/xdatabase/xredis"
//	"github.com/gin-gonic/gin"
//	"testing"
//	"time"
//)
//
//func TestRun(t *testing.T) {
//	xlog.WithDebugger()
//	xapi.AddResource(xdb.Init)
//	xapi.AddResource(xredis.Init)
//	xapi.AddResource(xqueue.NewClient)
//	engine := New(&WsConfig{
//		Auth: true,
//	})
//	engine.RegisterMsgHandler(1, func(c *Context) error {
//		xlog.Infof(c, "收到消息 %s", c.Msg.Message)
//		return nil
//	})
//	engine.StartStatistics()
//	xapi.StartApi(func(r *gin.Engine) error {
//		engine.RegisterToGin(r, "/ws")
//		return nil
//	}, "apis", "test")
//}
//
//func TestName(t *testing.T) {
//
//	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwNjExMDIyMzUsImlzcyI6Imp4Iiwic3ViIjoiMVx1MDAyNmFkbWluIn0.weujYJcSP1t7EdVJUU_UdHAybpIwi_2re640dpB-kmU 截止日期 2026-04-27 16:59:59
//	fmt.Println(ctx.GenerateToken(1, time.Hour*24*365*10, "admin"))
//}
