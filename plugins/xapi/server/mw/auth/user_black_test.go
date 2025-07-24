package auth

import (
	"testing"
	"time"
)

type cache struct {
	bs map[int64]struct{}
}

func (c *cache) Add(userId int64) error {
	c.bs[userId] = struct{}{}
	return nil
}

func (c *cache) Remove(userId int64) error {
	delete(c.bs, userId)
	return nil
}

func (c *cache) IsBlack(userId int64) bool {
	_, ok := c.bs[userId]
	return ok
}

func TestName(t *testing.T) {
	c := &cache{bs: make(map[int64]struct{})}
	c.Add(2)
	backCache = c
	AddUserBlack(1)
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(2))
	t.Log(IsUserBlack(2))
	t.Log(IsUserBlack(3))
	t.Log(IsUserBlack(3))
	time.Sleep(time.Second * 3)
	t.Log(IsUserBlack(2))
	t.Log(IsUserBlack(3))
	t.Log(IsUserBlack(3))
	t.Log(IsUserBlack(3))
	t.Log(IsUserBlack(3))
	t.Log(IsUserBlack(1))
	time.Sleep(time.Second * 3)
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
	t.Log(IsUserBlack(1))
}
