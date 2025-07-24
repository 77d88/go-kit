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
func (c *Client) RandomNum(ctx context.Context, workId uint16) int32 {
	key := IdGeneratorPrefix + fmt.Sprintf(":%d", workId)
	spop := c.SPop(ctx, key)
	if spop.Err() == nil {
		num, _ := spop.Int()
		return int32(num)
	}
	step := 1000
	incrby := c.IncrBy(ctx, key+":offset", int64(step))

	// 初始化种子 100000~110000
	random := rand.Intn(10000) + 100000
	if int64(step) == incrby.Val() {
		incrby = c.IncrBy(ctx, key+":offset", int64(step+random))
	}
	array := make([]string, step)
	// 存入set中
	for i := 0; i < step; i++ {
		array[i] = fmt.Sprintf("%d", incrby.Val()-int64(i))
	}
	_ = c.SAdd(ctx, key, array)
	return c.RandomNum(ctx, workId)
}
