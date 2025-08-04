package main

import (
	"flag"
	"fmt"
	"github.com/77d88/go-kit/cmd/xf/route"
	"github.com/77d88/go-kit/cmd/xf/util"
	"os"
)

func main() {
	path := flag.String("f", "./route.yml", "配置文件地址 默认为 ./route.yml")
	mode := flag.String("m", "1", "1. 路由生成 2.handler包装器生成")
	updateFilePath := flag.String("uf", "", "跟新文件的路径")
	flag.Parse()
	util.V.Set("mode", *mode)
	if util.V.GetInt("mode") == 1 {
		util.V.Set("path", *path)
		util.InitConfig(*path)
		//route.Build()
	}

	directory, _ := util.GetCurrentWorkingDirectory()

	if util.V.GetInt("mode") == 2 {

		if *updateFilePath != "" {
			err := route.UpdateRunFunc(*updateFilePath)
			if err != nil {
				panic(err)
			}
		} else {
			err := route.ScanAndGenerateRunFunctions(directory)
			if err != nil {
				panic(err)
			}
		}

	}

	fmt.Println("Current working directory:", directory)
	fmt.Println("Executable path:", os.Args)

}
