package redismq

import (
	"sync"
	"time"

	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/hibiken/asynq"
)

type Client struct {
	*asynq.Client
}

var client *Client
var once sync.Once

func NewX() *Client {
	c, err := x.Config[Config]("redis") // 使用默认的redis配置
	if err != nil {
		xlog.Errorf(nil, "xqueue redis config error: %v", err)
		return nil
	}
	// 创建一个客户端
	return New(c)
}

func New(config *Config) *Client {
	once.Do(func() {
		if config == nil {
			xlog.Errorf(nil, "xqueue redis config is nil")
			return
		}
		if config.Addr == "" {
			xlog.Errorf(nil, "xqueue redis config addr is empty")
			return
		}
		c := &Client{
			Client: asynq.NewClient(asynq.RedisClientOpt{
				Addr:     config.Addr,
				Password: config.Pass,
				DB:       config.Db,
			}),
		}
		client = c
	})

	return client
}

func Send(typename string, msg interface{}) error {
	return send(typename, msg)
}

func SendAt(typename string, msg interface{}, t time.Time) error {
	return send(typename, msg, asynq.ProcessAt(t))
}

func send(typename string, msg interface{}, options ...asynq.Option) error {
	json, err := xparse.ToJSON(msg)
	if err != nil {
		return err
	}
	options = append(options, asynq.Retention(24*time.Hour), asynq.MaxRetry(7)) // 保存1天方便排查 最多重试7次
	_, err = client.Enqueue(asynq.NewTask(typename, []byte(json)), options...)
	return err
}

func (c *Client) Dispose() error {
	err := c.Close()
	if err != nil {
		xlog.Errorf(nil, "redismq close error: %v", err)
	} else {
		xlog.Infof(nil, "redismq close successful")
	}
	return err
}
