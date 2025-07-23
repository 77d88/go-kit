package xid

import (
	"github.com/77d88/go-kit/basic/xparse"
	"testing"
	"time"
)

func TestNextId(t *testing.T) {
	t.Log(NextId())

	toTime, _ := xparse.ToTime("2018-08-08 08:08:08", time.DateTime)
	t.Log(toTime.UnixMilli())
}
