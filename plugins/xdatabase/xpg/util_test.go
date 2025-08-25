package xpg

import (
	"testing"
	"time"
)

func TestExObj(t *testing.T) {
	obj, _ := extractDBObj(&User{
		BaseModel: BaseModel{
			CreatedTime: time.Now(),
			DeletedTime: time.Now(),
			UpdatedTime: time.Now(),
		},
		UpdateUser:   0,
		Password:     "",
		Username:     "123",
		Nickname:     "",
		Avatar:       nil,
		Email:        "123",
		ReLoginDesc:  "123",
		_codes:       nil,
		_isCalcCodes: false,
	})
	for k, v := range obj {
		t.Log(k, v)
	}
}
func TestNamedExpr(t *testing.T) {
	t.Log(NamedExpr("update user set username = :username,name = :username where id = :id", map[string]interface{}{
		"username": "123",
		"id":       1,
	}))
}
