package main

import (
	"flag"
	"github.com/77d88/go-kit/cmd/route"
	"github.com/77d88/go-kit/cmd/util"
)

func main() {
	path := flag.String("f", "./route.yml", "配置文件地址 默认为 ./route.yml")
	mode := flag.String("m", "1", "1. 路由生成 2.handler包装器生成")
	flag.Parse()
	util.InitConfig(*path)
	util.V.Set("mode", *mode)
	util.V.Set("path", *path)
	if util.V.GetInt("mode") == 1 {
		route.Build()
	}

}
