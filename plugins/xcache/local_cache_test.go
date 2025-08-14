package xcache

import (
	"testing"
	"time"

	"github.com/77d88/go-kit/plugins/xlog"
)

type Res struct {
	ID          int
	CreatedTime time.Time
	UpdatedTime time.Time
	AliEtag     *string
}

func TestLocalCache_Set(t *testing.T) {
	fci00 := func() (interface{}, error) {
		println("123133")
		return &Res{
			ID: 2,
		}, nil
	}
	_, err := Once("test", time.Minute, fci00)
	_, err = Once("test", time.Minute, fci00)
	_, err = Once("test", time.Minute, fci00)
	_, err = Once("test2", time.Minute, fci00)
	res, err2 := Once("test", time.Minute, fci00)

	xlog.Errorf(nil, "err %v %v", err, nil)
	xlog.Errorf(nil, "err %v %v", err, nil)
	xlog.Errorf(nil, "err %+v %v", res, err2)
	//xlog.Errorf(nil, "err %v %v", err, nil == nil)

}

func funcName(key string) {

}

func funcName2(key string) {

}
