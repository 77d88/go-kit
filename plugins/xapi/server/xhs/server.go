package xhs

import (
	"context"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type ServerConfig struct {
	Port     int      `json:"port"`
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
	XE     *xe.Engine
	routes []string
}

func New(e *xe.Engine) *HttpServer {
	var cfg ServerConfig
	e.Cfg.ScanKey("server", &cfg)
	if cfg.Rate <= 0 {
		cfg.Rate = 100
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	loggerInit(&cfg)
	server := &HttpServer{Engine: engine, Config: &cfg, XE: e}
	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusOK, xerror.New("405", CodeMethodNotAllowed))
	})
	engine.NoRoute(func(c *gin.Context) {
		xlog.Debugf(newCtx(c, server).Copy(), "Not Found ")
		c.JSON(http.StatusOK, xerror.New("404", CodeNotFound))
	})
	engine.Use(WarpHandleMw(serverHandler))
	generatedDefaultRegister(server)
	e.MustProvide(func() *HttpServer {
		return server
	})
	return server
}

func (x *HttpServer) Start() {
	// 初始化http 服务
	x.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", x.Config.Port),
		Handler: x.Engine,
	}
	// 服务连接
	xlog.Infof(nil, "start success  prot: %d  production %v name: %v [%dms] [%droute]", x.Config.Port, !x.Config.Debug, x.Config.Name, x.XE.RunTime().Milliseconds(), len(x.routes))
	if err := x.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		xlog.Fatalf(nil, "listen: %s\n", err)
	}
}

func (x *HttpServer) Shutdown() {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := x.srv.Shutdown(timeout); err != nil {
		xlog.Errorf(nil, "http server stop error: %v", err)
	}
	xlog.Infof(nil, "http server stop success!!")
}

func (x *HttpServer) Use(fs HandlerMw) *HttpServer {
	x.Engine.Use(WarpHandleMw(fs))
	return x
}

func (x *HttpServer) POST(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return x.register(x.Engine.POST, "POST", path, handler, fs...)
}

func (x *HttpServer) GET(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return x.register(x.Engine.GET, "GET", path, handler, fs...)
}

func (x *HttpServer) PUT(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return x.register(x.Engine.PUT, "PUT", path, handler, fs...)
}

func (x *HttpServer) DELETE(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return x.register(x.Engine.DELETE, "DELETE", path, handler, fs...)
}

func (x *HttpServer) ANY(path string, handler Handler, fs ...HandlerMw) *HttpServer {
	return x.register(x.Engine.Any, "ANY", path, handler, fs...)
}

func (x *HttpServer) register(fc func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes, method string, path string, handler Handler, fs ...HandlerMw) *HttpServer {
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
	x.routes = append(x.routes, fmt.Sprintf("%s %s %s", method, path, location))
	fc(path, append(xarray.Map(fs, func(index int, item HandlerMw) gin.HandlerFunc {
		return WarpHandleMw(item)
	}), WarpHandle(handler))...)
	return x
}

func (x *HttpServer) Name() string {
	return fmt.Sprintf("http server  prot: %d  production %v name: %v ", x.Config.Port, !x.Config.Debug, x.Config.Name)
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

// WrapWithDI 依赖注入包装器 特殊情况使用 正常情况使用代码生成器来处理 尽量不使用反射调用
// 包装方法，将 *Ctx 作为第一个参数 其余参数使用容器注入，并返回 interface{}, error
//
//	示例 ： func a(c *xhs.Ctx, db *xdb.DataSource, redis *xredis.RedisClient) (interface{}, error) {
//		  如果容器中没有注册，则返回错误吗 8000 如果有 则正常生成路由函数
//		}
func (x *HttpServer) WrapWithDI(method any) func(*Ctx) (interface{}, error) {

	if method == nil {
		return func(ctx *Ctx) (interface{}, error) {
			return nil, xerror.New(CodeInvoke)
		}
	}
	if reflect.TypeOf(method).Kind() != reflect.Func {
		return func(ctx *Ctx) (interface{}, error) {
			return nil, xerror.New(CodeInvoke)
		}
	}

	methodValue := reflect.ValueOf(method)
	methodType := methodValue.Type()

	// 检查方法参数：第一个必须是 *Ctx，返回必须是 interface{}, error
	if methodType.NumIn() < 1 || methodType.In(0) != reflect.TypeOf((*Ctx)(nil)) ||
		methodType.NumOut() != 2 ||
		methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return func(ctx *Ctx) (interface{}, error) {
			return nil, xerror.New(CodeInvoke).SetMsg("method parameter error")
		}
	}

	// 构造一个符合 method 参数类型的函数签名用于注入
	paramTypes := make([]reflect.Type, methodType.NumIn()-1)
	for i := 1; i < methodType.NumIn(); i++ {
		paramTypes[i-1] = methodType.In(i)
	}

	inst, err := x.XE.GetInstValue(paramTypes...)
	if err != nil {
		xlog.Errorf(nil, "dig get inst error: %v", err)
		return func(ctx *Ctx) (interface{}, error) {
			return nil, xerror.New(CodeInvoke)
		}
	}
	// 返回包装后的处理函数
	return func(ctx *Ctx) (interface{}, error) {
		args := append([]reflect.Value{reflect.ValueOf(ctx)}, inst...)
		results := methodValue.Call(args)
		if len(results) != 2 {
			return nil, xerror.New(CodeInvoke)
		}
		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}
		return results[0].Interface(), nil
	}
}
