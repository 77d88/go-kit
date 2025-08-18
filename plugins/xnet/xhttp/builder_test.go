package xhttp

import (
	"net/http"
	"testing"

	"github.com/77d88/go-kit/plugins/xlog"
)

func TestName(t *testing.T) {

	xlog.WithDebugger()
	execute, err := NewPost("https://www.baidu.com").Execute()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(execute))

}

func TestCookies(t *testing.T) {
	cookie, err := http.ParseCookie("a=2;b=3")
	if err != nil {
		t.Error(err)
	}
	for i, i2 := range cookie {
		t.Log(i, i2.Name)

	}
}
