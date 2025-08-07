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
	fci00 := func() (interface{}, error) {
		println("123133")
		return &Res{
			ID: 1,
		}, nil
	}
	err := Once("test", nil, time.Minute, fci00)
	err = Once("test", nil, time.Minute, fci00)
	err = Once("test", nil, time.Minute, fci00)
	err = Once("test", nil, time.Minute, fci00)
	err = Once("test", nil, time.Minute, fci00)

	xlog.Errorf(nil, "err %v %v", err, nil)
	xlog.Errorf(nil, "err %v %v", err, nil)
	//xlog.Errorf(nil, "err %v %v", err, nil == nil)

}

func funcName(key string) {

}

func funcName2(key string) {

}
