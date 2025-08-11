package dbmw

import (
	xhs2 "github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
)

func TranManager() xhs2.HandlerMw {
	return func(c *xhs2.Ctx) {
		defer func() {

			if err := recover(); err != nil {
				tran := xdb.GetCtxTran(c)
				if tran == nil {
					return
				}
				if r := tran.DB.Rollback(); r.Error != nil {
					xlog.Errorf(c, "rollback error: %v", r.Error)
					panic(r.Error)
				}
				xlog.Debugf(c, "request error db rollback")
				panic(err)
			}
		}()
		c.Next() // 先执行业务逻辑

		tran := xdb.GetCtxTran(c)
		if tran == nil {
			return
		}

		if c.Errors.Last() != nil {
			if r := tran.DB.Rollback(); r.Error != nil {
				xlog.Errorf(c, "rollback error: %v", r.Error)
				return
			}
			xlog.Debugf(c, "request error db rollback %v", c.Errors.String())
			return
		}

		if r := tran.DB.Commit(); r.Error != nil {
			xlog.Errorf(c, "commit error: %v", r.Error)
			// 交给回滚 处理
			if r := tran.DB.Rollback(); r.Error != nil {
				xlog.Errorf(c, "rollback error: %v", r.Error)
			}
		}
	}
}
