package xcache

import (
	"github.com/77d88/go-kit/plugins/xlog"
	"testing"
	"time"
)

type Res struct {
	ID          int
	CreatedTime time.Time
	UpdatedTime time.Time
	AliEtag     *string
}

func TestLocalCache_Set(t *testing.T) {
	var res Res
	fci00 := func() (interface{}, error) {
		println("123133")
		return &Res{
			ID: 1,
		}, nil
	}
	err := Once("test", &res, time.Minute, fci00)
	err = Once("test", &res, time.Minute, fci00)
	err = Once("test", &res, time.Minute, fci00)
	err = Once("test", &res, time.Minute, fci00)
	var res2 Res
	err = Once("test", &res2, time.Minute, fci00)

	xlog.Errorf(nil, "err %v %v", err, res)
	xlog.Errorf(nil, "err %v %v", err, res2)
	xlog.Errorf(nil, "err %v %v", err, res2 == res)

}

func funcName(key string) {

}

func funcName2(key string) {

}
