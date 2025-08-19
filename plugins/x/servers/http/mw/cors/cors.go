package cors

import (
	"net/http"
	"strings"

	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
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

func New() xhs.HandlerMw {
	cs := x.ConfigStringSlice("server.cors")
	oriCors := []string{
		"://localhost",
		"null",
		"://127.0.0.1",
		"://192.168",
		"i-ii.top",
	}
	c := &CorsConfig{
		origins: append(oriCors, cs...),
	}
	Cors = c
	return corsMw(c)
}

// corsMw 处理跨域请求,支持options访问
func corsMw(config *CorsConfig) xhs.HandlerMw {
	return func(c *xhs.Ctx) {
		// Access-Control-Allow-Credentials=true和Access-Control-Allow-Origin="*"有冲突
		// 故Access-Control-Allow-Origin需要指定具体得跨域origin
		method := c.Request.Method               // 请求方法
		origin := c.Request.Header.Get("Origin") // 请求头部
		if method == http.MethodOptions {
			if !config.CorsCheck(origin) { // 没有支持的域 拒绝访问
				c.Send(xerror.New("no support cors origin"))
				c.Abort()
				return
			}
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Credentials", "true") //  跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Origin", origin)      // 这是允许访问所有域
			requestMethods := c.Request.Header.Get("Access-Control-Request-Method")
			c.Header("Access-Control-Allow-Methods", requestMethods) // 服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			requestHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
			c.Header("Access-Control-Allow-Headers", requestHeaders) // 所有头部
			c.Header("Access-Control-Max-Age", "172800")             // 缓存请求信息 单位为秒
		}

		// 放行所有OPTIONS方法
		if method == http.MethodOptions {
			// 确保响应头被写入
			c.AbortWithStatus(http.StatusNoContent)
			return
			//return
		}
		c.Next()
	}
}
