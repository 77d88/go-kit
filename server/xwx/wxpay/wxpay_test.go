package wxpay

//import (
//	"context"
//	"github.com/77d88/go-kit/basic/xconfig/redis_scanner"
//	"github.com/77d88/go-kit/basic/xid"
//	"testing"
//)
//
//func Test_Run(t *testing.T) {
//	InitWith(redis_scanner.Default("default"))
//	order := PlaceOrder{
//		Amount: 1,
//		Openid: "oy7lF4xm3Gj4P0VfxQcvTnAuhjjg",
//		PaySn:  xid.NextId(),
//		Desc:   "test",
//	}
//
//	jsapi, err := TransactionJsapiPrepayId(context.TODO(), &order)
//	if err != nil {
//		t.Error(err)
//	} else {
//		t.Log(jsapi)
//	}
//}
