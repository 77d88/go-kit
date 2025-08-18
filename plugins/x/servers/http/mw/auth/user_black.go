package auth

import (
	"context"
	"time"

	"github.com/77d88/go-kit/plugins/xlog"
)

type UserBlackCache interface {
	Add(userId int64) error
	Remove(userId int64) error
	IsBlack(userId int64) bool
}
type ms struct {
	cacheTime int64 // unix 时间戳秒级
	is        bool
}

// 用户黑名单
var userBlacklist = make(map[int64]*ms)

var backCache UserBlackCache // 黑名单缓存

func SetCache(cache UserBlackCache) {
	backCache = cache
}

func AddUserBlack(userId int64) {
	if backCache != nil {
		err := backCache.Add(userId)
		if err != nil {
			xlog.Errorf(context.TODO(), "add user black error:%s", err.Error())
		}
	}
	userBlacklist[userId] = &ms{
		cacheTime: time.Now().Unix(),
		is:        true,
	}
}
func RemoveUserBlack(userId int64) {
	if backCache != nil {
		err := backCache.Remove(userId)
		if err != nil {
			xlog.Errorf(context.TODO(), "remove user black error:%s", err.Error())
		}
	}
	delete(userBlacklist, userId)
}
func IsUserBlack(userId int64) bool {
	v, ok := userBlacklist[userId]
	if ok { // 本地存在数据存在的情况
		if backCache == nil {
			return v.is // 全靠本地数据咯
		}
		rate := int64(60) // 存在 每分钟同步一次缓存
		if !v.is {        // 不存在 30秒同步一次
			rate = 30
		}
		if time.Now().Unix()-v.cacheTime > rate {
			delete(userBlacklist, userId) // 删除本地 下次就走缓存加载咯
		}
		return v.is
	} else {
		if backCache == nil { // 没有就是没有
			return false
		} else {
			black := backCache.IsBlack(userId)
			userBlacklist[userId] = &ms{ // 缓存同步
				cacheTime: time.Now().Unix(),
				is:        black,
			}
			return black
		}
	}

}
