package route

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/77d88/go-kit/basic/xstr"
	"github.com/77d88/go-kit/cmd/xf/util"
)

// RouteConfig 定义路由配置结构
type RouteConfig struct {
	WorkPath string
	Biz      string
	Routes   []Route `mapstructure:"routes"`
}

// Route 定义单个路由配置
type Route struct {
	Path     string `mapstructure:"path"`
	Module   string
	Methods  []string  `mapstructure:"methods"`
	Handlers []Handler `mapstructure:"handlers"`
}

// Handler 定义处理器配置
type Handler struct {
	Name    string   `mapstructure:"name"`
	Module  string   `mapstructure:"module"`
	Route   string   `mapstructure:"route"`
	Remark  string   `mapstructure:"remark"`
	Auth    bool     `mapstructure:"auth"`
	Methods []string `mapstructure:"methods,omitempty"`
}

var biz = "biz"

func GenRouteAll(b string) {

	var config RouteConfig
	if err := util.V.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		os.Exit(1)
	}
	if b != "" {
		config.Biz = b
	}
	if config.Biz == "" {
		config.Biz = biz
	}
	biz = config.Biz
	// 获取当前工作目录
	wd, err := util.GetCurrentWorkingDirectory()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// 创建 biz 目录
	bizDir := filepath.Join(wd, biz)
	if err := os.MkdirAll(bizDir, 0755); err != nil {
		fmt.Printf("Error creating %s directory: %v\n", biz, err)
		os.Exit(1)
	}

	//生成代码
	if err := generateCode(config, bizDir); err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Code generated successfully in /%s directory!\n", biz)
}

// HandlerInfo 包含处理程序及其路由信息
type HandlerInfo struct {
	BasePath   string
	Route      string
	Methods    []string
	Handler    Handler
	Remark     string
	ModuleName string
	BizModel   string
	FullPath   string
}

func generateCode(config RouteConfig, bizDir string) error {
	handlers := make([]HandlerInfo, 0)
	// 为每个处理器生成代码
	for _, route := range config.Routes {
		currBizDir := filepath.Join(bizDir, route.Module)
		for i, handler := range route.Handlers {
			methods := handler.Methods
			if len(methods) == 0 {
				methods = route.Methods
			}
			if handler.Name == "" {
				fmt.Printf("%s handler[%d] name is empty \n", route.Module, i)
				continue
			}

			if handler.Route == "" {
				handler.Route = "/" + handler.Name
			}
			if handler.Module == "" {
				handler.Module = handler.Name
			}
			if handler.Remark == "" {
				handler.Remark = xstr.Capitalize(handler.Name)
			}

			handlerInfo := HandlerInfo{
				BasePath:   route.Path,
				Route:      handler.Route,
				Methods:    methods,
				Handler:    handler,
				Remark:     handler.Remark,
				ModuleName: handler.Module,
				BizModel:   route.Module,
				FullPath:   fmt.Sprintf("%s%s", route.Path, handler.Route),
			}
			handlers = append(handlers, handlerInfo)

			// 为每个处理器生成代码，moduleName作为子目录
			if err := generateModuleCode(handler.Module, handlerInfo, currBizDir); err != nil {
				return err
			}
		}
	}
	// 在bizDir 生成所有路由的register
	if err := genRegister(config, handlers, bizDir); err != nil {
		return err
	}

	return nil
}

func generateModuleCode(moduleName string, handler HandlerInfo, bizDir string) error {
	// 创建模块目录
	moduleDir := filepath.Join(bizDir, moduleName)
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return err
	}
	handlerFile := filepath.Join(moduleDir, "handler.go")
	// 如果文件存在跳过
	if _, err := os.Stat(handlerFile); !os.IsNotExist(err) {
		return nil
	}

	// 生成 route.go 文件
	if err := generateRouteFile(moduleName, handler, moduleDir); err != nil {
		return err
	}

	// 生成 run.go 文件
	if err := UpdateRunFunc(filepath.Join(moduleDir, "handler.go")); err != nil {
		return err
	}

	return nil
}

func generateRouteFile(moduleName string, handler HandlerInfo, moduleDir string) error {
	tmpl := `package {{.ModuleName}}

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
{{- if .Handler.Auth}}
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
{{- end}}
)

// {{.Handler.Remark}}
type request struct {
}

func handler(c *xhs.Ctx, r *request) (interface{},error) {
	return nil,nil
}

func run(c *xhs.Ctx) (interface{}, error) {
	if r,err := xhs.ShouldBind[request](c);err != nil {
		return nil, xerror.New("参数错误").SetCode(xhs.CodeParamError).SetInfo("参数错误: %+v", err)
	}else{
		return handler(c, &r)
	}
}

func Register(xsh *xhs.HttpServer) {
{{- range .Methods}}
	xsh.{{.}}("{{$.BasePath}}{{$.Route}}", run{{if $.Handler.Auth}},auth.ForceAuth {{end}})
{{- end}}
}
`

	data := struct {
		ModuleName string
		BasePath   string
		Route      string
		Methods    []string
		Handler    Handler
	}{
		ModuleName: moduleName,
		BasePath:   handler.BasePath,
		Route:      handler.Route,
		Methods:    handler.Methods,
		Handler:    handler.Handler,
	}

	t, err := template.New("handler").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("error formatting handler code: %v\nCode:\n%s", err, buf.String())
	}

	// 写入文件
	filename := filepath.Join(moduleDir, "handler.go")
	return os.WriteFile(filename, formatted, 0644)
}

func genRegister(config RouteConfig, modules []HandlerInfo, bizDir string) error {
	tmpl := `// Code generated by xf. DO NOT EDIT.

package {{.Biz}}

import (
{{- range .Modules}}
	{{.Alias}} "{{.ImportPath}}"
{{- end}}
	"github.com/77d88/go-kit/plugins/x/servers/http/xhs"
)

func Register(xsh *xhs.HttpServer) {
{{- range .Modules}}
	{{.Alias}}.Register(xsh)
{{- end}}		
}
`

	// 准备模板数据
	type moduleTemplateData struct {
		ImportPath  string
		PackageName string
		Alias       string
		FullPath    string
	}

	var moduleData []moduleTemplateData
	processedModules := make(map[string]bool) // 避免重复导入

	for _, module := range modules {
		importPath := fmt.Sprintf("%s/%s/%s/%s", config.WorkPath, biz, module.BizModel, module.ModuleName)
		if !processedModules[importPath] {
			// 构造导入路径，假设是相对于当前模块的路径
			packageName := strings.ReplaceAll(module.ModuleName, "/", "_")

			moduleData = append(moduleData, moduleTemplateData{
				ImportPath:  importPath,
				Alias:       xstr.CamelCase(module.BizModel + "_" + packageName),
				PackageName: packageName,
				FullPath:    fmt.Sprintf("%s%s", module.BasePath, module.Route),
			})

			processedModules[module.ModuleName] = true
		}
	}

	data := struct {
		Modules []moduleTemplateData
		Biz     string
	}{
		Modules: moduleData,
		Biz:     biz,
	}

	t, err := template.New("register").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("error formatting register code: %v\nCode:\n%s", err, buf.String())
	}

	// 写入文件
	filename := filepath.Join(bizDir, "register.go")
	return os.WriteFile(filename, formatted, 0644)
}
