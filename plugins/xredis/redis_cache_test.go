package xredis

// import (
//	"context"
//	"fmt"
//	"testing"
//	"time"
//	"xcore/xcache"
//	"xcore/xcv"
// )
//
// type Res struct {
//	Key          int
//	CreatedTime time.Time
//	UpdatedTime time.Time
//	AliEtag     *string
// }
//
// func TestLocalCache_Set(t *testing.T) {
//	Init(Config{
//		Addr: "127.0.0.1:6379",
//		Pass: "",
//		Db:   0,
//	})
//	xcache.UseCache(NewRedisCache("text", Cmd))
//
//	k1 := "tx"
//	_ = xcache.Set(context.Background(), k1, "123", time.Second*1)
//	funcName(k1)
//
//	time.Sleep(time.Second * 2)
//	funcName(k1)
//
//	xcache.Del(context.Background(), k1)
//	funcName(k1)
//
//	k2 := "tx2"
//	_ = xcache.Set(context.Background(), k2, Res{
//		Key:          3,
//		CreatedTime: time.Time{},
//		UpdatedTime: time.Time{},
//		AliEtag:     xcv.V2p("3312"),
//	}, time.Second*2)
//	funcName(k2)
//
//	xcache.Del(context.Background(), k2)
//	funcName(k2)
//	xcache.Setf(context.Background(), k2, time.Second*20, func(ctx context.DefaultAppContext) (interface{}, error) {
//		return Res{
//			Key:          4,
//			CreatedTime: time.Now(),
//			UpdatedTime: time.Now(),
//			AliEtag:     xcv.V2p("3312"),
//		}, nil
//	})
//	funcName(k2)
//
//	get, ok := xcache.Get(context.Background(), k2)
//	if ok {
//		res, err := xcv.JsonToBean(get.(string), &Res{})
//		if err != nil {
//			t.Error(err)
//		}
//		fmt.Printf("%#v\n", res)
//	}
//
// }
//
// func funcName(key string) {
//	get, err := xcache.Get(context.Background(), key)
//	if err {
//		b, ok := get.(string)
//		if ok {
//			println("有数据")
//			println(b)
//		} else {
//			println("有数据 但转换错误")
//		}
//	} else {
//		println("没有数据")
//	}
// }
