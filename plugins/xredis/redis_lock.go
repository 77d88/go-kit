package xredis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
	"log"
	"sync"
	"time"
)

type RedisLock struct {
	client    *RedisClient
	key       string
	value     string // 唯一标识锁持有者
	ttl       time.Duration
	cancelCtx context.Context
	cancelFn  context.CancelFunc
	mutex     sync.Mutex
}

var (
	ErrLockNotAcquired = errors.New("redis lock: not acquired")
	ErrLockNotHeld     = errors.New("redis lock: not held by this instance")
)

func NewLock(key string, ttl time.Duration, name ...string) *RedisLock {
	db, err := Get(name...)
	if err != nil {
		xlog.Errorf(nil, "数据库[%s]链接不存在", name)
		return nil
	}
	return &RedisLock{
		client: db,
		key:    key,
		ttl:    ttl,
	}
}

// NewLockByClient 创建Redis锁实例 自定义客户端
func NewLockByClient(client *RedisClient, key string, ttl time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}
}

// Lock 获取锁（带重试机制）
func (l *RedisLock) Lock(ctx context.Context, retryInterval time.Duration, maxRetry int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 生成唯一标识（防止误删其他客户端的锁）
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return fmt.Errorf("generate token failed: %w", err)
	}
	l.value = base64.StdEncoding.EncodeToString(token)

	var retryCount int
	for {
		// 尝试获取锁
		acquired, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
		if err != nil {
			return fmt.Errorf("redis setnx failed: %w", err)
		}

		if acquired {
			// 启动自动续期
			l.cancelCtx, l.cancelFn = context.WithCancel(context.Background())
			go l.startWatchdog()
			return nil
		}

		// 重试逻辑
		if retryCount >= maxRetry {
			return ErrLockNotAcquired
		}
		retryCount++

		select {
		case <-time.After(retryInterval):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// 自动续期（看门狗机制）
func (l *RedisLock) startWatchdog() {
	ticker := time.NewTicker(l.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 延长锁过期时间
			l.mutex.Lock()
			ok, err := l.client.Expire(l.cancelCtx, l.key, l.ttl).Result()
			l.mutex.Unlock()

			if err != nil || !ok {
				log.Printf("watchdog renew failed: %v", err)
				return
			}
		case <-l.cancelCtx.Done():
			return
		}
	}
}

// Unlock 释放锁（Lua脚本保证原子性）
func (l *RedisLock) Unlock() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cancelFn != nil {
		l.cancelFn() // 停止看门狗
	}

	// Lua脚本确保只有锁持有者才能释放
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end`

	_, err := l.client.Eval(context.Background(), script, []string{l.key}, l.value).Result()
	if err != nil {
		return fmt.Errorf("unlock script failed: %w", err)
	}
	return nil
}

// Run 获取锁并运行函数
func (l *RedisLock) Run(retryInterval time.Duration, maxRetry int, f func() error) error {
	if err := l.Lock(context.Background(), retryInterval, maxRetry); err != nil {
		return err
	}
	defer l.Unlock()

	return f()
}
