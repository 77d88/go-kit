package main

import (
	"flag"
	"github.com/77d88/go-kit/v2/cmd/route"
	"github.com/77d88/go-kit/v2/cmd/util"
)

func main() {
	path := flag.String("f", "./route.yml", "配置文件地址 默认为 ./route.yml")
	flag.Parse()
	util.InitConfig(*path)
	route.Build()
}
