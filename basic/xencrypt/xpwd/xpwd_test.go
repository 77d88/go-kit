package xpwd

import "testing"

func TestName(t *testing.T) {

	password, err := HashPassword("123456")
	if err != nil {
		t.Error(err)
	}
	t.Log(password)
	t.Log(CheckPasswordHash("123456", password))
}
