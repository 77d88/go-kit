package initDataRows

import (
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
)

// 重置初始化数据
type response struct {
}

type request struct {
	DbName    string   `json:"dbName"`
	TableName []string `json:"tableName"`
	Restore   bool     `json:"restore"` // 是否重新建立
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if len(r.DbName) == 0 {
		r.DbName = xdb.DefaultDbLinkStr
	}

	for k, dist := range xdb.RegisterModels {
		if k == r.DbName {
			for t, v := range dist {
				if v.GetID() <= 0 {
					continue
				}
				// 创建一个v实例
				d := xcore.NewBy(v)
				result := xdb.C(c, r.DbName).First(&d, v.GetID())
				if xdb.IsNotFound(result.Error) { // 记录不存在，创建新记录
					result := xdb.C(c, r.DbName).Create(v)
					if result.Error != nil {
						return nil, result.Error
					}
					xlog.Infof(nil, "创建记录【%s】[%s]成功: %v", t, k, v.GetID())
				} else {
					// 删除在创建
					if r.Restore { // 重建建立所有数据
						if result := xdb.C(c, r.DbName).Unscoped().Unscoped().Delete(v); result.Error != nil {
							return nil, result.Error
						}
						if result := xdb.C(c, r.DbName).Create(v); result.Error != nil {
							return nil, result.Error
						}
						xlog.Infof(nil, "记录【%s】[%s]=> %v 已存删除重新建立成功", t, k, v.GetID())
					} else {
						xlog.Infof(nil, "记录【%s】[%s]=> %v 已存在不重新建立", t, k, v.GetID())
					}
				}
			}

		}
	}
	return
}

func Register(xsh *xhs.HttpServer) {
	xsh.POST("/pro/db/initDataRows", run(), auth.ForceAuth)
}
