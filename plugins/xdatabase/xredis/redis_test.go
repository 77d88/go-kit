package xredis

import (
	"context"
	"fmt"
	"testing"

	"github.com/77d88/go-kit/plugins/xlog"
)

func TestGetUserNum(t *testing.T) {
	New(&Config{
		Addr: "127.0.0.1:6379",
		Pass: "jerry123!",
		Db:   0,
	})
	New(&Config{
		Addr:       "127.0.0.1:6379",
		Pass:       "jerry123!",
		Db:         1,
		DbLinkName: "test",
	})
	db, err := Get()
	if err != nil {
		xlog.Errorf(nil, "redis init error ,%v", err)
		t.Error(err)
		return
	}
	d2, _ := Get("test")
	for i := 0; i < 1; i++ {
		fmt.Printf("id -> %d \n", db.RandomNum(context.Background(), 2))
		fmt.Printf("id2 -> %d \n", d2.RandomNum(context.Background(), 1))
	}
}

//
//func TestConcurrentLimit(t *testing.T) {
//	with := New(&Config{
//		Addr: "127.0.0.1:6379",
//		Pass: "jerry123!",
//		Db:   0,
//	})
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		time.Sleep(1 * time.Second)
//		go func(id int) {
//			defer wg.Done()
//			err := with.LimitRun(context.Background(), "test", 1*time.Second, func() error {
//				fmt.Printf("run == > %v\n", id)
//				return nil
//			})
//			if err != nil {
//				return
//			}
//		}(i)
//	}
//	wg.Wait()
//}
//
//func TestConcurrentLock(t *testing.T) {
//	New(&Config{
//		Addr: "127.0.0.1:6379",
//		Pass: "jerry123!",
//		Db:   0,
//	})
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func(id int) {
//			defer wg.Done()
//			lock := NewLock("lockKey", 5*time.Second)
//			err := lock.Run(1*time.Second, 3, func() error {
//				fmt.Printf("run == > %v\n", id)
//				time.Sleep(1 * time.Second)
//				return nil
//			})
//			if err != nil {
//				fmt.Printf("err == > %d: %v\n", id, err)
//			}
//
//		}(i)
//	}
//	wg.Wait()
//}
//
//func Test_str(t *testing.T) {
//	ctx := context.Background()
//	with := New(&Config{
//		Addr: "127.0.0.1:6379",
//		Pass: "jerry123!",
//		Db:   0,
//	})
//	t.Run("get", func(t *testing.T) {
//		array := make([]int, 200)
//		// 存入set中
//		for i := 0; i < 200; i++ {
//			array[i] = i
//		}
//		with.SAdd(ctx, "abcxxx", array)
//		//Rds.Sadd("abcxxx", strings)
//		str := with.SPop(ctx, "abcxxx")
//		fmt.Printf("xxx == > %s\n", str)
//		fmt.Printf("xxx == > %v\n", str.Err())
//	})
//
//	t.Run("get", func(t *testing.T) {
//		str := with.Get(ctx, "abcxxx")
//		fmt.Printf("xxx == > %s\n", str.String())
//		fmt.Printf("xxx == > %v\n", str.Err())
//	})
//
//	t.Run("get", func(t *testing.T) {
//		str := with.Get(ctx, "abc")
//		fmt.Printf("%s\n", str.String())
//	})
//	t.Run("set", func(t *testing.T) {
//		_ = with.Set(ctx, "abc", "123", -1)
//	})
//	t.Run("set", func(t *testing.T) {
//		_ = with.SetEx(ctx, "abc", "123", 60*time.Hour)
//	})
//}
