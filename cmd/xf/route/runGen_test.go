package route

import (
	"github.com/77d88/go-kit/cmd/xf/util"
	"testing"
)

func TestName(t *testing.T) {

	directory, _ := util.GetCurrentWorkingDirectory()
	err := ScanAndGenerateRunFunctions(directory)

	if err != nil {
		panic(err)
	}
}

func TestUpdate(t *testing.T) {
	err := UpdateRunFunc("G:\\development\\project\\AAAAtools\\go\\commonv2\\examples\\xhsserver\\biz\\user_service\\v2\\create\\handler.go")
	if err != nil {
		panic(err)
	}
}

func TestGenAll(t *testing.T)  {
	util.InitConfig("G:\\development\\project\\AAAAtools\\go\\commonv2\\cmd\\xf\\route.yml")
	GenRouteAll("test")
	
}
