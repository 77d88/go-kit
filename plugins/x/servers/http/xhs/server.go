package xhs

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port     string   `json:"port"`
	Rate     int      `json:"rate"`
	Name     string   `json:"name"`
	Debug    bool     `json:"debug"`
	Cors     []string `json:"cors"`
	LogLevel int      `json:"logLevel"` // -1 trace 0 debug 1 info 2 warn 3 error 4 fatal  debug=false的时候使用json
}

type HttpServer struct {
	Engine *gin.Engine
	srv    *http.Server
	Config *ServerConfig
	fcs    []func()
	mws    []HandlerMw
}

func New() *HttpServer {
	c, err := x.Config[ServerConfig]("server")
	if err != nil {
		xlog.Fatalf(nil, "config error: %v", err)
	}
	if c.Rate <= 0 {
		c.Rate = 100
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	loggerInit(c)
	server := &HttpServer{Engine: engine, Config: c}
	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusOK, xerror.New("405", CodeMethodNotAllowed))
	})
	engine.NoRoute(func(c *gin.Context) {
		xlog.Debugf(newCtx(c, server).Copy(), "Not Found ")
		c.JSON(http.StatusOK, xerror.New("404", CodeNotFound))
	})
	engine.Use(WarpHandleMw(serverHandler))
	generatedDefaultRegister(server)
	x.Use(func() *HttpServer {
		return server
	})
	if err != nil {
		xlog.Fatalf(nil, "provide error: %v", err)
	}
	return server
}

func (h *HttpServer) Start() {
	// 初始化路由
	c1 := make(chan struct{})
	// 初始化基础mw
	go func() {
		for _, mw := range h.mws {
			h.Engine.Use(WarpHandleMw(mw))
		}
		c1 <- struct{}{}
	}()
	for _, fc := range h.fcs {
		fc()
	}
	<-c1
	// 初始化http 服务
	h.srv = &http.Server{
		Addr:    net.JoinHostPort("", h.Config.Port),
		Handler: h.Engine,
	}
	// 服务连接
	xlog.Infof(nil, "start success  prot: %s  production %v name: %v [%dms] [%droute]", h.Config.Port, !h.Config.Debug, h.Config.Name, time.Since(x.Info().StartTime).Milliseconds(), len(h.fcs))
	if err := h.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		xlog.Fatalf(nil, "listen: %s\n", err)
	}
}

func (h *HttpServer) Shutdown() {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.srv.Shutdown(timeout); err != nil {
		xlog.Errorf(nil, "http server stop error: %v", err)
	}
	xlog.Infof(nil, "http server stop success!!")
}

func (h *HttpServer) Use(fs HandlerMw) *HttpServer {
	h.mws = append(h.mws, fs)
	return h
}

func (h *HttpServer) POST(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return h.register(h.Engine.POST, "POST", path, handler, fs...)
}

func (h *HttpServer) GET(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return h.register(h.Engine.GET, "GET", path, handler, fs...)
}

func (h *HttpServer) PUT(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return h.register(h.Engine.PUT, "PUT", path, handler, fs...)
}

func (h *HttpServer) DELETE(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return h.register(h.Engine.DELETE, "DELETE", path, handler, fs...)
}

func (h *HttpServer) ANY(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return h.register(h.Engine.Any, "ANY", path, handler, fs...)
}

func (h *HttpServer) register(fc func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes, method string, path string, handler Handler, fs ...HandlerMw) *HttpServer {
	h.fcs = append(h.fcs, func() {
		// 获取 handler 的文件名和行号
		handlerPtr := runtime.FuncForPC(reflect.ValueOf(handler).Pointer())
		var location string
		if handlerPtr != nil {
			file, line := handlerPtr.FileLine(handlerPtr.Entry())
			location = fmt.Sprintf("%s:%d", file, line)
		} else {
			location = "unknown"
		}
		xlog.Debugf(nil, "register %s %s %s mw[%d] ", method, path, location, len(fs))
		fc(path, append(xarray.Map(fs, func(index int, item HandlerMw) gin.HandlerFunc {
			return WarpHandleMw(item)
		}), WarpHandle(handler))...)
	})
	return h
}

func (h *HttpServer) Name() string {
	return fmt.Sprintf("http server  prot: %s  production %v name: %v ", h.Config.Port, !h.Config.Debug, h.Config.Name)
}

func loggerInit(cfg *ServerConfig) {
	if cfg.Debug {
		xlog.WithDebugger()
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			xlog.Tracef(nil, "endpoint %v %v %v %v", httpMethod, absolutePath, handlerName, nuHandlers)
		}
	} else {
		xlog.WithRelease()
	}
}
