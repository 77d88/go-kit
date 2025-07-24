package xtime

import "time"

// In 是否在时间区间内
func In(check, start, end time.Time) bool {
	return check.After(start) && check.Before(end)
}

// BeginOfMinute  返回t的分钟开始时间.
func BeginOfMinute(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 0, 0, t.Location())
}

// EndOfMinute  返回t的分钟结束时间.
func EndOfMinute(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), t.Minute(), 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfHour 返回t的小时开始时间.
func BeginOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}

// EndOfHour 返回t的小时结束时间.
func EndOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfDay 返回t的当日开始时间.
func BeginOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回t的当日结束时间.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// BeginOfWeek 返回从周开始，默认周从星期日开始.
func BeginOfWeek(t time.Time, beginFrom time.Weekday) time.Time {
	y, m, d := t.AddDate(0, 0, int(beginFrom-t.Weekday())).Date()
	beginOfWeek := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
	if beginOfWeek.After(t) {
		return beginOfWeek.AddDate(0, 0, -7)
	}
	return beginOfWeek
}

// EndOfWeek 返回结束周时间，默认周末为星期六.
func EndOfWeek(t time.Time, endWith time.Weekday) time.Time {
	y, m, d := t.AddDate(0, 0, int(endWith-t.Weekday())).Date()
	var endWithWeek = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	if endWithWeek.Before(t) {
		endWithWeek = endWithWeek.AddDate(0, 0, 7)
	}
	return endWithWeek
}

// BeginOfMonth 返回月初时间	.
func BeginOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 返回月底时间。
func EndOfMonth(t time.Time) time.Time {
	return BeginOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// BeginOfYear 返回年初的日期时间。
func BeginOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 返回年底的日期时间。
func EndOfYear(t time.Time) time.Time {
	return BeginOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// IsLeapYear 检查 param year 是否为闰年。
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DayOfYear 返回参数 date 't' 在一年中的哪一天。
func DayOfYear(t time.Time) int {
	y, m, d := t.Date()
	firstDay := time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
	nowDate := time.Date(y, m, d, 0, 0, 0, 0, t.Location())

	return int(nowDate.Sub(firstDay).Hours() / 24)
}

// BetweenDays returns the number of days between two times.
func BetweenDays(start, end time.Time) int {
	duration := end.Sub(start)
	days := int(duration.Hours() / 24)

	return days
}

// BetweenSeconds 返回两次之间的秒数。
func BetweenSeconds(t1 time.Time, t2 time.Time) int64 {
	index := t2.Unix() - t1.Unix()
	return index
}

// Min 返回给定时间中最早的时间.
func Min(t1 time.Time, times ...time.Time) time.Time {
	minTime := t1

	for _, t := range times {
		if t.Before(minTime) {
			minTime = t
		}
	}

	return minTime
}

// Max 返回给定时间中的最新时间.
func Max(t1 time.Time, times ...time.Time) time.Time {
	maxTime := t1

	for _, t := range times {
		if t.After(maxTime) {
			maxTime = t
		}
	}

	return maxTime
}

// MaxMin 返回给定时间中的最新时间和最早时间.
func MaxMin(t1 time.Time, times ...time.Time) (maxTime time.Time, minTime time.Time) {
	maxTime = t1
	minTime = t1

	for _, t := range times {
		if t.Before(minTime) {
			minTime = t
		}

		if t.After(maxTime) {
			maxTime = t
		}
	}

	return maxTime, minTime
}
