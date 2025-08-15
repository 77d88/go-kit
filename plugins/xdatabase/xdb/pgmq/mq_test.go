package pgmq

import (
	"testing"

	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
)

func TestName(t *testing.T) {

	xdb.New(&xdb.Config{
		Dns:    xdb.FastDsn("127.0.0.1", 5432, "postgres", "jerry123!", "zyv2"),
		Logger: true,
	})

	mq := New(&Config{})
	mq.RegisterHandler(MsgType(0), func(msg *Queue) (bool, error) {
		t.Logf("收到消息,%v", msg)
		return true, nil
	})
	//err := mq.Send("yesy")
	//go func() {
	//	for {
	//		message, err2 := mq.SendMessage(context.Background(), Msg{
	//			Msg:   "12333",
	//			Retry: 0,
	//		})
	//		if err2 != nil {
	//			t.Error(err2)
	//		}
	//		t.Log(message)
	//		time.Sleep(time.Second * 2)
	//	}
	//
	//}()
	//mq.Start()
	err := mq.Start()
	if err != nil {
		t.Error(err)
		panic(err)
	}

	select {}

}
