package xhs

import (
	"github.com/77d88/go-kit/plugins/xe"
)

func generatedDefaultRegister(r *HttpServer) {

	// 获取服务状态
	r.GET("/x/sys/info/status", func(c *Ctx) (interface{}, error) {
		return map[string]interface {
		}{
			"name":      r.Config.Name,                 // 服务名称
			"runTime":   xe.E.RunTime().Milliseconds(), // 运行时长
			"startTime": xe.E.Info.StartTime,
		}, nil
	})
	r.GET("/ping", func(c *Ctx) (interface{}, error) {
		return "pong", nil
	})

}
