package xctx

import (
	"github.com/77d88/go-kit/plugins/xapi/server/xhs"
	"testing"
)

func TestName(t *testing.T) {

	c := xhs.NewTestContext()

	val := SetVal(c, "key", "value")

	t.Log(val.Value("key"))
}
