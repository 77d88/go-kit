package redis_scanner

import (
	"os"
	"testing"
)

func TestName(t *testing.T) {
	err := os.Setenv("V_CONFIG_REDIS_ADDR", "redis://default:yzztest@127.0.0.1:6666/0")
	if err != nil {
		panic(err)
	}
	Default("base", "mini", "test")
}
