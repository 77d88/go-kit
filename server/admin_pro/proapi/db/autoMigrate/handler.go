package autoMigrate

import (
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 自动同步数据库结构
type response struct {
}

type request struct {
	DbName    string   `json:"dbName"`
	TableName []string `json:"tableName"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if len(r.DbName) == 0 {
		r.DbName = xdb.DefaultDbLinkStr
	}
	db, err := xdb.GetDB(r.DbName)
	if err != nil {
		return
	}
	i := make([]interface{}, 0, len(xdb.RegisterModels))

	for k, v := range xdb.RegisterModels {
		if r.DbName == k {
			for t, tv := range v {
				if len(r.TableName) == 0 || xarray.Contain(r.TableName, t) {
					i = append(i, tv)
				}
			}
		}
	}
	err = db.DB.AutoMigrate(i...)

	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth, pro.HansPermission(pro.Per_SuperAdmin))
}
