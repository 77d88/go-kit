package route

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xredis"
)

const (
	ForceAuthMwPath = "github.com/77d88/go-kit/plugins/xapi/server/mw/auth" // 强制认证中间件包名
	DataSourcePath  = "github.com/77d88/go-kit/plugins/xdb"                 // 数据源包名
	RedisDataPath   = "github.com/77d88/go-kit/plugins/redis"               // redis数据源包名
)

type response struct {
}

type request struct {
}

//go:generate ../xf -m=2
func handler(c *xhs.Ctx, r *request, db *xdb.DB, redis *xredis.Client) (interface{}, error) {
	return response{}, nil
}
