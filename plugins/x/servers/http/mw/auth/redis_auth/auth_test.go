package redis_auth

import (
	"github.com/77d88/go-kit/plugins/x/servers/http/mw/auth"
	"github.com/77d88/go-kit/plugins/xredis"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	a2 := &Auth{
		Client: xredis.Init(&xredis.Config{
			Addr: "127.0.0.1:6379",
			Db:   0,
			Pass: "jerry123!",
		}),
		Prefix:      "test",
		AutoRenewal: true,
	}

	login, err := a2.Login(1, auth.WithMaxLoginNum(2))
	if err != nil {
		t.Error(err)
	}
	t.Log(login)
	time.Sleep(time.Second * 2)

}

func TestLogout(t *testing.T) {

	a2 := &Auth{
		Client: xredis.Init(&xredis.Config{
			Addr: "127.0.0.1:6379",
			Db:   0,
			Pass: "jerry123!",
		}),
		Prefix:      "test",
		AutoRenewal: true,
	}
	// MTo5MDQ1MTg2OTY1NTQ1NjU
	err := a2.Logout("MTo5MDQ1MTg2OTY1NTQ1NjU")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)

}

func Test22(t *testing.T) {
	t.Log((time.Hour*24*7).Seconds())
	t.Log(time.Duration(604800000000000).Seconds())
}
