package xcache

import "fmt"

type Value struct {
	Val interface{}
}

type Warper interface {
	WarpKey(key string) string
	WarpValue(val interface{}) *Value
	UnWarpValue(val *Value) (interface{}, bool)
}

type PrefixWarp struct {
	Prefix string
}

func (l *PrefixWarp) WarpKey(key string) string {
	return fmt.Sprintf("%s:%s", l.Prefix, key)
}

func (l *PrefixWarp) WarpValue(val interface{}) *Value {
	return &Value{
		Val: val,
	}
}
func (l *PrefixWarp) UnWarpValue(val *Value) (interface{}, bool) {
	if val == nil {
		return nil, false
	}
	return val.Val, true
}
