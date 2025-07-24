package xtime

import "time"

type TimeInterval struct {
	start map[string]time.Time
}

func NewTimeInterval() *TimeInterval {
	m := make(map[string]time.Time)
	m["default"] = time.Now()
	return &TimeInterval{
		start: m,
	}
}

func (i TimeInterval) AddGroup(group string) *TimeInterval {
	i.start[group] = time.Now()
	return &i
}

// IntervalMs 间隔毫秒
func (i TimeInterval) IntervalMs() int64 {
	t := i.start["default"]
	return time.Now().Sub(t).Milliseconds()
}

// IntervalS 间隔秒
func (i TimeInterval) IntervalS() float64 {
	t := i.start["default"]
	return time.Now().Sub(t).Seconds()
}

// IntervalM 间隔分钟
func (i TimeInterval) IntervalM() float64 {
	t := i.start["default"]
	return time.Now().Sub(t).Minutes()
}

// IntervalGroupMs 分组间隔毫秒
func (i TimeInterval) IntervalGroupMs(group string) int64 {
	t := i.start[group]
	return time.Now().Sub(t).Milliseconds()
}

// IntervalGroupS 分组间隔秒
func (i TimeInterval) IntervalGroupS(group string) float64 {
	t := i.start[group]
	return time.Now().Sub(t).Seconds()
}

// IntervalGroupM 分组间隔分钟
func (i TimeInterval) IntervalGroupM(group string) float64 {
	t := i.start[group]
	return time.Now().Sub(t).Minutes()
}
