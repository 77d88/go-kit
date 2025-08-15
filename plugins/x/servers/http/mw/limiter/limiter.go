package limiter

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"golang.org/x/time/rate"
)

// New 全局的限流器 令牌桶限流
func New() xhs.HandlerMw {
	limit := x.ConfigInt("server.rate")
	if limit > 0 {
		limit = 100
	}
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
