package route

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xlog"
)

const (
	ForceAuthMwPath = "github.com/77d88/go-kit/plugins/xapi/server/mw/auth" // 强制认证中间件包名
	DataSourcePath  = "github.com/77d88/go-kit/plugins/xdb"                 // 数据源包名
	RedisDataPath   = "github.com/77d88/go-kit/plugins/redis"               // redis数据源包名
)

func Build() {

}

type response struct {
}

type request struct {
}

func handler(c *xhs.Ctx, r request) (interface{}, error) {
	return response{}, nil
}

func Run() xhs.Handler {
	return func(c *xhs.Ctx) (r interface{}, e error) {
		// 事务可以再这里开启
		return warpHandler(c)
	}
}

// 这个是一个包裹器,是自动生成的代码
func warpHandler(c *xhs.Ctx) (r interface{}, e error) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	return handler(c, r)
}
