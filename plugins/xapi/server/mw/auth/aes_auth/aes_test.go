package aes_auth

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {

	a := New()

	login, err := a.Login(1, "admin")
	if err != nil {
		t.Error(err)
	}
	t.Log(login)
	data := a.VerificationToken(login.Token)
	if data.Err != nil {
		t.Error(data.Err)
	}
	t.Log(data)
	t.Log(data.ExpireTime.Unix())

	t.Log(generateToken(1, time.Hour*24*365*10, defaultAuthV2Key))

}
