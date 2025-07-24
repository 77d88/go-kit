package xhs

import (
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestName(t *testing.T) {

	cfg := ServerConfig{
		Port:     9981,
		Rate:     100,
		Name:     "xxx",
		Debug:    true,
		LogLevel: 0,
	}
	engine := gin.New()
	//loggerInit(&cfg)
	server := &HttpServer{Engine: engine, Config: &cfg, QuitSignal: make(chan os.Signal)}
	server.GET("/ping", func(ctx *Ctx) (interface{}, error) {
		return nil, nil
	})
	server.Start()

}
