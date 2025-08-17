package xredis

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrLimitNotAcquired = errors.New("limit: not acquired")
)

// LimitRun  锁的变种 只能自动释放
func (c *client) LimitRun(ctx context.Context, key string, ttl time.Duration, f func() error) error {
	acquired, err := c.SetNX(ctx, key, "1", ttl).Result()
	// 尝试获取限流器
	if err != nil {
		return fmt.Errorf("redis setnx failed: %w", err)
	}
	if acquired {
		return f()
	} else {
		return ErrLimitNotAcquired
	}
}
