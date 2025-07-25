package dbmw

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
)

func TranManager() xhs.Handler {
	return func(c *xhs.Ctx) (interface{}, error) {
		defer func() {

			if err := recover(); err != nil {
				tran := xdb.GetCtxTran(c)
				if tran == nil {
					return
				}
				if r := tran.DB.Rollback(); r.Error != nil {
					xlog.Errorf(c, "rollback error: %v", r.Error)
				}
				xlog.Debugf(c, "request error db rollback")
			}
		}()
		c.Next() // 先执行业务逻辑

		tran := xdb.GetCtxTran(c)
		if tran == nil {
			return nil, nil
		}

		if c.Errors.Last() != nil {
			if r := tran.DB.Rollback(); r.Error != nil {
				xlog.Errorf(c, "rollback error: %v", r.Error)
				return nil, r.Error
			}
			xlog.Debugf(c, "request error db rollback %v", c.Errors.String())
			return nil, nil
		}

		if r := tran.DB.Commit(); r.Error != nil {
			xlog.Errorf(c, "commit error: %v", r.Error)
			panic(r.Error) // 交给回滚 处理
		}

		return nil, nil
	}
}
