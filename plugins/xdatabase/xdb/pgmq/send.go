package pgmq

import (
	"context"
	"encoding/json"
	"time"
)

type SendMessageOptions struct {
	Delay int // 延迟收到秒数
	Type  MsgType
	Retry int // 重试几次 1重试一次
}

// SendMessageOption 是用于配置 SendMessageOptions 的函数类型。
type SendMessageOption func(*SendMessageOptions)

// WithDelay 是一个 SendMessageOption，用于设置消息的延迟可见性。
func WithDelay(seconds int) SendMessageOption {
	return func(opts *SendMessageOptions) {
		opts.Delay = seconds
	}
}
func WithType(t MsgType) SendMessageOption {
	return func(opts *SendMessageOptions) {
		opts.Type = t
	}
}
func WithRetry(retry int) SendMessageOption {
	return func(opts *SendMessageOptions) {
		opts.Retry = retry
	}
}

// Send 发送消息到队列
func (xq *XQueue) Send(message interface{}, opts ...SendMessageOption) error {
	return xq.SendWithDelay(message, opts...)
}

// SendWithDelay 发送延迟消息到队列
func (xq *XQueue) SendWithDelay(message interface{}, opts ...SendMessageOption) error {
	var opt SendMessageOptions
	for _, o := range opts {
		o(&opt)
	}
	if opt.Retry < 0 {
		opt.Retry = 0
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	result := xq.db.WithCtx(context.Background()).Create(&Queue{
		Message:      string(msgBytes),
		State:        QueueStatePending,
		Num:          0,
		Retry:        int16(opt.Retry + 1), // 这个是最终执行几次 最小为1 设置为0 则不执行
		Type:         int(opt.Type),
		DeliveryTime: time.Now().Add(time.Duration(opt.Delay) * time.Second),
	})
	return result.Error
}
