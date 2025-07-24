package xerror

import (
	"fmt"
)

const (
	ErrorCodeSystem = -1 // 系统错误
)

type XError interface {
	XError() error
}

type Error struct {
	Msg  string `json:"msg,omitempty"`
	Code int    `json:"code"`
	Info string `json:"info,omitempty"`
}

func (b *Error) XError() error {
	return b
}

// 实现 error 接口的 Error 方法
func (b *Error) Error() string {
	return fmt.Sprintf("Error %d: %s", b.Code, b.Msg)
}

// New creates a new Error instance
// msg: Can be Error, *Error, string, *string, int, or any other type
// code: Optional error code (default -1)
func New(msg interface{}, code ...int) *Error {
	c := ErrorCodeSystem
	if len(code) > 0 {
		c = code[0]
	}

	switch v := msg.(type) {
	case Error:
		return &Error{Msg: v.Msg, Code: v.Code, Info: v.Info}
	case *Error:
		return &Error{Msg: v.Msg, Code: v.Code, Info: v.Info}
	case string:
		return &Error{Msg: v, Code: c}
	case *string:
		msg := "System error"
		if v != nil {
			msg = *v
		}
		return &Error{Msg: msg, Code: c}
	case int:
		return &Error{Msg: "System error", Code: v}
	default:
		return &Error{
			Msg:  "System error",
			Info: fmt.Sprintf("Error: %v", msg),
			Code: c,
		}
	}
}

func (b *Error) SetMsg(msg string, args ...interface{}) *Error {
	b.Msg = fmt.Sprintf(msg, args...)
	return b
}
func (b *Error) SetCode(code int) *Error {
	b.Code = code
	return b
}

func (b *Error) SetInfo(info string, args ...interface{}) *Error {
	b.Info = fmt.Sprintf(info, args...)
	return b
}

func Newf(msg string, args ...interface{}) *Error {
	return &Error{
		Msg:  fmt.Sprintf(msg, args...),
		Code: ErrorCodeSystem,
	}
}

// IsXError 检查错误是否为特定类型(Error或其指针类型)
// 参数:
//
//	err: 要检查的错误对象
//
// 返回值:
//
//	true 当err是Error类型或其指针类型时
//	false 其他情况
//
// 注意:
//
//	使用errors.As()进行错误类型匹配是更符合Go习惯的做法
func IsXError(err interface{}) bool {
	switch err.(type) {
	case *Error, Error:
		return true
	default:
		return false
	}
}
