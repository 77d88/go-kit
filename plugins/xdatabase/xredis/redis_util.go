package xredis

import (
	"context"
	"fmt"
	"math/rand"
)

const (
	IdGeneratorPrefix = "ID_GENERATOR"
)

// RandomNum 获取随机数唯一数
func RandomNum(ctx context.Context, workId uint16, name ...string) (int, error) {
	script, err := GetScript(name...)
	if err != nil {
		return 0, err
	}
	key := IdGeneratorPrefix + fmt.Sprintf(":%d", workId)
	// 执行Lua脚本
	//return script.Eval(ctx, lua_randomNum, []string{key}, 1000, rand.Intn(10000)+100000).Int()
	cmd := script.EvalSha(ctx, script_randomNum, []string{key}, 1000, rand.Intn(10000)+100000)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Int()
}
