package cors

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"net/http"
	"strings"
)

type CorsConfig struct {
	origins []string
}

var (
	Cors *CorsConfig
)

func (c *CorsConfig) CorsCheck(origin string) bool {
	if origin == "" {
		return false
	}
	_, ok := xarray.Find(c.origins, func(i int, item string) bool {
		return strings.Contains(origin, item)
	})
	return ok
}

func Check(origin string) bool {
	if Cors == nil {
		return true
	}
	return Cors.CorsCheck(origin)
}

func New(cfg *xhs.ServerConfig) xhs.Handler {
	oriCors := []string{
		"://localhost",
		"null",
		"://127.0.0.1",
		"://192.168",
		"i-ii.top",
	}
	c := &CorsConfig{
		origins: append(oriCors, cfg.Cors...),
	}
	Cors = c
	return corsMw(c)
}

// corsMw 处理跨域请求,支持options访问
func corsMw(config *CorsConfig) xhs.Handler {
	return func(c *xhs.Ctx) (interface{}, error) {
		// Access-Control-Allow-Credentials=true和Access-Control-Allow-Origin="*"有冲突
		// 故Access-Control-Allow-Origin需要指定具体得跨域origin
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部
		if method == http.MethodOptions {
			if !config.CorsCheck(origin) { // 没有支持的域 拒绝访问
				return nil, xerror.New("no support cors origin")
			}
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Credentials", "true")                               //  跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Origin", origin)                                    // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization,Device,Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma") // 所有头部
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                                                                                                                                                          // 缓存请求信息 单位为秒
		}

		// 放行所有OPTIONS方法
		if method == http.MethodOptions {
			c.JSON(http.StatusOK, "")
			c.Abort()
			return nil, nil
		}
		return nil, nil
	}
}
