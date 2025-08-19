package xredis

import (
	"context"
	"errors"
	"sync"

	"github.com/redis/go-redis/v9"
)

type EvalCmd struct {
	err error
	*redis.Cmd
}

func (c *EvalCmd) Err() error {
	if c.err != nil {
		return c.err
	}
	if c.Cmd.Err() != nil {
		return c.Cmd.Err()
	}
	return nil
}

// ScriptManager 脚本管理器
type ScriptManager struct {
	client  *redis.Client
	scripts map[string]string
	shas    map[string]string
	mu      sync.RWMutex
}

// NewScriptManager 创建脚本管理器
func NewScriptManager(client *redis.Client) *ScriptManager {
	return &ScriptManager{
		client:  client,
		scripts: make(map[string]string),
		shas:    make(map[string]string),
	}
}

// LoadScript 加载脚本
func (sm *ScriptManager) LoadScript(ctx context.Context, name, script string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sha, err := sm.client.ScriptLoad(ctx, script).Result()
	if err != nil {
		return err
	}

	sm.scripts[name] = script
	sm.shas[name] = sha
	return nil
}

// EvalSha 执行脚本
func (sm *ScriptManager) EvalSha(ctx context.Context, name string, keys []string, args ...interface{}) *EvalCmd {
	sm.mu.RLock()
	sha, exists := sm.shas[name]
	sm.mu.RUnlock()

	if !exists {
		s, e := lua_map[name] // 默认脚本
		if e {
			err := sm.LoadScript(ctx, name, s)
			if err != nil {
				return &EvalCmd{err: err, Cmd: redis.NewCmd(ctx)}
			} else {
				s, exists := sm.shas[name]
				if !exists {
					return &EvalCmd{err: errors.New(name + " script not exists"), Cmd: redis.NewCmd(ctx)}
				}
				sha = s
			}
		} else {
			return &EvalCmd{err: errors.New(name + " script not exists"), Cmd: redis.NewCmd(ctx)}
		}

	}

	cmd := sm.client.EvalSha(ctx, sha, keys, args...)

	// 如果脚本不存在，重新加载
	if cmd.Err() != nil && cmd.Err().Error() == "NOSCRIPT No matching script. Please use EVAL." {
		sm.mu.Lock()
		// 双重检查
		if newSha, exists := sm.shas[name]; exists && newSha != sha {
			// 其他goroutine已经重新加载了脚本
			sm.mu.Unlock()
			return &EvalCmd{Cmd: sm.client.EvalSha(ctx, newSha, keys, args...)}
		}

		// 重新加载脚本
		if script, exists := sm.scripts[name]; exists {
			newSha, err := sm.client.ScriptLoad(ctx, script).Result()
			if err == nil {
				sm.shas[name] = newSha
				sha = newSha
			}
		}
		sm.mu.Unlock()

		if sha != "" {
			cmd = sm.client.EvalSha(ctx, sha, keys, args...)
		}
	}

	return &EvalCmd{Cmd: cmd}
}
