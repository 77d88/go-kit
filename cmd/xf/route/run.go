package route

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xe"
	"github.com/77d88/go-kit/plugins/xredis"
)

func Run() xhs.Handler {
	var db *xdb.DB
	var redis *xredis.Client
	xe.E.MustInvoke(func(p1 *xdb.DB, p2 *xredis.Client) {
		db = p1
		redis = p2
	})
	return func(c *xhs.Ctx) (interface{}, error) {
		r := request{}
		err := c.ShouldBind(&r)
		if err != nil {
			return nil, xerror.New("参数错误").SetCode(xhs.CodeParamError).SetInfo("参数错误: %+v", err)
		}
		return handler(c, &r, db, redis)
	}
}
