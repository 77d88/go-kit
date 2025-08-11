package dbModels

import (
	"github.com/77d88/go-kit/basic/xmap"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdb"
)

// 数据库表列表
type response struct {
	DbName    string   `json:"dbName"`
	TableName []string `json:"tableName"`
}

type request struct {
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	models := xdb.RegisterModels

	result := make([]response, 0, len(models))

	for k, v := range models {
		result = append(result, response{
			DbName:    k,
			TableName: xmap.Keys(v),
		})
	}

	return result, nil
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
