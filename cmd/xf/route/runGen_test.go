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
