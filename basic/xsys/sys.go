package xsys

import (
	"bytes"
	"fmt"
	"github.com/77d88/go-kit/basic/xcore"
	"github.com/77d88/go-kit/basic/xparse"
	"github.com/77d88/go-kit/basic/xtype"
	"os"
	"os/exec"
	"runtime"
)

// StackTrace returns 堆栈的字符串
func StackTrace(all bool) string {
	var buf bytes.Buffer

	// Reserve a larger initial buffer to reduce reallocations
	initialSize := 5120
	buf.Grow(initialSize)

	for {
		// Allocate a temporary buffer with the current buffer size
		tempBuf := make([]byte, buf.Cap())

		// Capture the stack trace
		size := runtime.Stack(tempBuf, all)

		// Check for errors (e.g., memory allocation failure)
		if size == 0 {
			// Handle the error case, e.g., xlog an error or return an empty string
			return ""
		}

		// Write the captured stack trace to the buffer
		buf.Write(tempBuf[:size])

		// If the captured stack trace fits in the buffer, break
		if size < len(tempBuf) {
			break
		}

		// Double the buffer size to accommodate the stack trace
		buf.Grow(initialSize)
	}

	// Return the captured stack trace as a string
	return buf.String()
}

// OsEnvGet 获取环境变量，如果为空则返回默认值
func OsEnvGet(key string, defaultValue string) string {
	env := os.Getenv(key)
	if xcore.IsZero(env) {
		return defaultValue
	}
	return env
}

func OsEnvGetNumber[T xtype.Numer](key string, defaultValue T) T {
	env := OsEnvGet(key, "")
	if env == "" {
		return defaultValue
	}
	if f, e := xparse.ToNumber[T](env); e != nil {
		return defaultValue
	} else {
		return f
	}
}

func Restart() {
	fmt.Println("强制重启程序")
	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return
	}

	// 使用 exec.Command 重新启动当前程序
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 设置环境变量等
	cmd.Env = os.Environ()

	// 启动新进程
	err = cmd.Start()
	if err != nil {
		fmt.Println("启动新进程失败:", err)
		return
	}

	// 等待新进程启动
	fmt.Println("新进程已启动，PID:", cmd.Process.Pid)

	// 退出当前进程
	os.Exit(0)
}
