package xqueue

import (
	"context"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/basic/xconfig/redis_scanner"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/hibiken/asynq"
	"testing"
	"time"
)

func init() {
	xlog.WithDebugger()
	xconfig.Init(redis_scanner.NewEnv(), "apis", "test")
	NewClient()
}

func TestName(t *testing.T) {

	//for {
	C.Enqueue(asynq.NewTask("test", []byte("xxxxx")))
	C.Enqueue(asynq.NewTask("test2", []byte("2222222"), asynq.ProcessAt(time.Now().Add(time.Second*10)), asynq.Retention(2*time.Minute)))
	time.Sleep(time.Millisecond)
	//}

}

func TestConsumer(t *testing.T) {
	// Use asynq.HandlerFunc adapter for a handler function
	NewServer(asynq.Config{
		Queues: map[string]int{
			"default": 1,
			"high":    6,
			"slow":    4,
		},
	}, map[string]ConsumerFunc{
		"test": func(ctx context.Context, payload string) error {
			xlog.Infof(nil, "payload test: %v", payload)
			return nil
		},
		"test2": func(ctx context.Context, payload string) error {
			xlog.Infof(nil, "payload test2: %v", payload)
			return nil
		},
	})

	select {}
}
