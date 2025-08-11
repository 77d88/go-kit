package limiter

import (
	"github.com/77d88/go-kit/basic/xerror"
	xhs2 "github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"golang.org/x/time/rate"
)

// Limiter 全局的限流器 令牌桶限流
func Limiter(limit int) xhs2.HandlerMw {
	limiter := rate.NewLimiter(rate.Limit(limit), limit*2)
	return func(c *xhs2.Ctx) {
		if !limiter.Allow() {
			c.Send(xerror.New("frequent!!").SetCode(xhs2.CodeCurrentLimiting))
			c.Abort()
			return
		}
		c.Next()
	}
}
