package redismq

import (
	"context"
	"testing"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/basic/xid"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/hibiken/asynq"
)

var config *Config

func init() {
	xlog.WithDebugger()
	config = &Config{
		Addr: "127.0.0.1:6666",
		Pass: "token1",
	}
}

func TestName(t *testing.T) {
	New(config)
	//for {
		_ = send("test", "default_"+xid.NextIdStr(),asynq.Queue("default"))
		//_ = send("test2", "high_"+xid.NextIdStr(),asynq.Queue("high"))
		//_ = send("test3", "slow_"+xid.NextIdStr(),asynq.Queue("slow"))
		//time.Sleep(time.Millisecond * 10)
	//}

}

func TestConsumer(t *testing.T) {
	// Use asynq.HandlerFunc adapter for a handler function
	NewServer(config, asynq.Config{
		Queues: map[string]int{
			"default": 1,
			"high":    6,
			"slow":    4,
		},
	}, map[string]ConsumerFunc{
		"test": func(ctx context.Context, payload string) error {
			xlog.Infof(nil, "payload test: %v", payload)
			return xerror.New("test error")
		},
		"test2": func(ctx context.Context, payload string) error {
			xlog.Infof(nil, "payload test2: %v", payload)
			return nil
		},
		"test3": func(ctx context.Context, payload string) error {
			xlog.Infof(nil, "payload test3: %v", payload)
			return nil
		},
	})

	select {}
}
