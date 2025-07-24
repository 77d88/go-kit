package xlog

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	//WithDebugger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxLogParam, map[string]interface{}{
		"testx": "test",
	})
	DefaultLogger.Info().Fields(map[string]interface{}{
		"test": "test",
	}).Msg("123")
	Tracef(ctx, "3333")

	split := strings.Split("C:/Users/Administrator/go/pkg/mod/codeup.aliyun.com/64812ff4d2963c5649be2afc/go-xcore@v1.1.7-0.20250408082323-3708d31a58b7/xapi/cors.go：111", "go-xcore")
	a := strings.Split(split[1], "/")
	a[0] = strings.Split(a[0], "-")[0]
	var as string
	for i, v := range a {
		if i == 0 {
			as += v
		} else {
			as += "/" + v
		}
	}
	t.Log(a)
	fmt.Printf("go-xcore%s", as)

	//G:/development/project/yzz/v3/apis/biz/service/orderservice/order.go:237

}
