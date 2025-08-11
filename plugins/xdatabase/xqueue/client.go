package xqueue

import (
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/hibiken/asynq"
	"time"
)

type Client struct {
	*asynq.Client
}

var C *Client

func NewClient() *Client {
	// 创建一个客户端
	c := &Client{
		Client: asynq.NewClient(getRedisConfig()),
	}
	C = c
	return c
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
	options = append(options, asynq.Retention(10*time.Minute)) // 默认保存十分钟方便排查罢了
	_, err = C.Enqueue(asynq.NewTask(typename, []byte(json)), options...)
	return err
}

func (c *Client) Dispose() error {
	return c.Close()
}
