package xredis

import (
	"context"
	"fmt"
	"github.com/77d88/go-kit/plugins/xlog"
)

type UserBlackCache struct {
	Prefix string
	client *Client
}

func NewUserBlackCache(prefix string, name ...string) *UserBlackCache {
	get, err := Get(name...)
	if err != nil {
		return nil
	}
	return &UserBlackCache{
		Prefix: prefix,
		client: get,
	}
}

func (u UserBlackCache) Add(userId int64) error {
	add := u.client.SAdd(context.TODO(), "black:"+u.Prefix, fmt.Sprintf("%d", userId))
	if add.Err() != nil {
		xlog.Errorf(context.TODO(), "add user black error:%s", add.Err().Error())
		return add.Err()
	}
	return nil
}

func (u UserBlackCache) Remove(userId int64) error {
	remove := u.client.SRem(context.TODO(), "black:"+u.Prefix, fmt.Sprintf("%d", userId))
	if remove.Err() != nil {
		xlog.Errorf(context.TODO(), "remove user black error:%s", remove.Err())
		return remove.Err()
	}
	return nil

}

func (u UserBlackCache) IsBlack(userId int64) bool {
	member := u.client.SIsMember(context.TODO(), "black:"+u.Prefix, fmt.Sprintf("%d", userId))
	if member.Err() != nil {
		xlog.Errorf(context.TODO(), "is black error:%s", member.Err().Error())
		return false
	}
	return member.Val()

}
