package api_menu

import "github.com/77d88/go-kit/plugins/xdb"

type apiRequest struct {
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
