package xconfig

import (
	"fmt"
	"time"
)

func clog(level string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s][%s] %s\n", timestamp, level, msg)
	fmt.Printf(logMsg, args...)
}

func InfoLog(msg string, args ...interface{}) {
	clog("INFO", msg, args...)
}

func ErrorLog(msg string, args ...interface{}) {
	clog("ERROR", msg, args...)
}

func WarnLog(msg string, args ...interface{}) {
	clog("WARN", msg, args...)
}
