package limiter

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"golang.org/x/time/rate"
)

// Limiter 全局的限流器 令牌桶限流
func Limiter(limit int) xhs.HandlerMw {
	limiter := rate.NewLimiter(rate.Limit(limit), limit*2)
	return func(c *xhs.Ctx) {
		if !limiter.Allow() {
			c.Send(xerror.New("frequent!!").SetCode(xhs.CodeCurrentLimiting))
			c.Abort()
			return
		}
		c.Next()
	}
}
