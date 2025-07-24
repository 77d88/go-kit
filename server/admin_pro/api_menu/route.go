package api_menu

import (
	"github.com/77d88/go-kit/basic/xarray"
	pro2 "github.com/77d88/go-kit/server/admin_pro/pro"
	"sort"
)

type Router struct {
	Path      string    `json:"path"`
	Component string    `json:"component,omitempty"`
	Redirect  string    `json:"redirect,omitempty"`
	Meta      *Meta     `json:"meta,omitempty"`
	Children  []*Router `json:"children,omitempty"`
	Name      string    `json:"name,omitempty"`
}

type Meta struct {
	ID         int64  `json:"id,string"`
	Title      string `json:"title"`
	Icon       string `json:"icon,omitempty"`
	HideInMenu bool   `json:"hideInMenu"`
	KeepAlive  bool   `json:"keepAlive"`
	Sort       int    `json:"sort"`
	NoLevel    bool   `json:"noLevel"`
	Params     string `json:"params,omitempty"`
}

func ConvertMenusToRouter(allMenus []*pro2.Menu, roles ...string) []*Router {
	menuMap := make(map[int64]*pro2.Menu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}

	rootMenus := make([]*Router, 0, 200)
	for _, menu := range allMenus {
		if !hasParent(menu.ID, allMenus) && hasPermission(menu, roles) {
			if router := buildRouterTree(menu, menuMap, roles); router != nil {
				rootMenus = append(rootMenus, router)
			}
		}
	}

	// 根菜单排序
	sort.Slice(rootMenus, func(i, j int) bool {
		return rootMenus[i].Meta.Sort < rootMenus[j].Meta.Sort
	})

	return rootMenus
}

func buildRouterTree(menu *pro2.Menu, menuMap map[int64]*pro2.Menu, roles []string) *Router {
	if !hasPermission(menu, roles) {
		return nil
	}
	router := &Router{
		Path:      menu.Path,
		Component: menu.ComponentPath,
		Redirect:  menu.Redirect,
		Name:      menu.Name,
		Meta: &Meta{
			ID:         menu.ID,
			Title:      menu.NameZh,
			Icon:       menu.MetaIcon,
			HideInMenu: menu.MetaHide,
			Sort:       menu.Sort,
			NoLevel:    menu.MetaNoLevel,
			Params:     menu.RouteParams,
		},
	}

	var children []*Router
	for _, childID := range menu.Children.ToSlice() {
		if childMenu, exists := menuMap[childID]; exists {
			if childRouter := buildRouterTree(childMenu, menuMap, roles); childRouter != nil {
				children = append(children, childRouter)
			}
		}
	}

	// 子菜单排序
	sort.Slice(children, func(i, j int) bool {
		return children[i].Meta.Sort < children[j].Meta.Sort
	})

	// 单子菜单自动提升逻辑
	if len(children) == 1 && !menu.MetaNoLevel {
		child := children[0]
		router.Path = child.Path
		router.Component = child.Component
		router.Name = child.Name
		router.Meta = child.Meta
		router.Children = child.Children
	} else {
		router.Children = children
	}

	// 如果处理后没有子菜单且是隐藏菜单，则返回nil
	// && router.Meta.HideInMenu 隐藏的菜单依然返回前端会隐藏 路由需要这个
	//if len(router.Children) == 0 {
	//	return nil
	//}

	return router
}

func hasParent(menuID int64, allMenus []*pro2.Menu) bool {
	for _, menu := range allMenus {
		for _, childID := range menu.Children.ToSlice() {
			if childID == menuID {
				return true
			}
		}
	}
	return false
}

func hasPermission(menu *pro2.Menu, userRoles []string) bool {
	if menu.Permission.IsEmpty() {
		return true
	}
	if xarray.Contain(userRoles, pro2.RoleSuperAdmin) {
		return true
	}

	for _, role := range userRoles {
		if menu.Permission.Contain(role) { // 任意包含
			return true
		}
	}
	return false
}
