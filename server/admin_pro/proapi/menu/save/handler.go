package save

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/server/admin_pro/pro"
)

// 菜单保存
type response struct {
}

type request struct {
	Id            int64          `json:"id"`
	Path          string         `json:"path"`
	ComponentPath string         `json:"componentPath"`
	Redirect      string         `json:"redirect"`
	Name          string         `json:"name"`
	NameZh        string         `json:"nameZh"`
	MataTitle     string         `json:"mataTitle"`
	MataKeywords  string         `json:"mataKeywords"`
	MetaIcon      string         `json:"metaIcon"`
	MetaHide      bool           `json:"metaHide"`
	Sort          int            `json:"sort"`
	MetaNoLevel   bool           `json:"metaNoLevel"`
	RouteParams   string         `json:"routeParams"`
	Children      *xdb.Int8Array `json:"children"`
}

func handler(c *xhs.Ctx, r *request) (resp interface{}, err error) {
	if r.Name == "" {
		return nil, xerror.New("参数错误:名称不能为空")
	}
	if result := xdb.Ctx(c).SaveMap(&pro.Menu{}, r, map[string]interface{}{
		"update_user": c.GetUserId(),
	}); result.Error != nil {
		return nil, result.Error
	}
	return
}

func Register(path string, xsh *xhs.HttpServer) {
	xsh.POST(path, run(), auth.ForceAuth)
}
