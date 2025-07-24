package xhs

import (
	"encoding/json"
	"fmt"
	"github.com/77d88/go-kit/basic/xarray"
	"github.com/77d88/go-kit/basic/xerror"
)

const DefaultErrorMsg = "系统错误"

func (c *Ctx) SendError(err interface{}) {
	e := xerror.New(err)
	c.Result = e
	_ = c.Error(e)
}

// Fatalf fatalConfig 统一快速处理错误
// 示例：
// 1. c.Fatalf("any2", xerror.newCtx("123")) 最后一个参数是xError 不打印错误日志 并在返回msg中附带错误信息
// 2. c.Fatalf(xerror.newCtx("123"), "用户端信息", "后台日志信息%s %v", "ss",1) 不打印日志 返回自定义信息 并打印日志(xlog.error(c,"后台日志信息%s %v", "ss",1))
// 3. c.Fatalf(xerror.newCtx("123"))  不打印日志 并在返回msg中附带错误信息(msg:123 code:-1)
// 4. c.Fatalf(xerror.newCtx("123"), xapi.FatalWithCode(100))  不打印日志 并在返回msg中附带错误信息 同时设置错误码 这种情况需要使用xapi.FatalWithMsg自定义返回自定义信息
// 5. c.Fatalf("123")    不打印错误日志 并在返回msg中附带错误信息(msg:123 code:-1)
// 6. c.Fatalf(errors.newCtx("123")  打印错误日志 并在返回info中附带错误信息 (msg:系统错误 code:-1 info:error 123)
// 7. c.Fatalf(errors.newCtx("123"), "错误", "执行")  执行打印错误信息  并在msg中附带自定义消息 (msg:错误 code:-1 info:error 123)
// 8. c.Fatalf("any","any2", errors.newCtx("123")) 最后一个参数是error 则打印错误日志 并在返回info中附带错误信息 (msg:系统错误 code:-1 info:error 123)
// 9. c.Fatalf(true, "用户端信息", "执行") 执行打印错误信息  并在msg中附带消息 (msg:用户端信息 code:-1 )
// 10. c.Fatalf(true, "错误", "hh %s", "ss", xapi.FatalWithCode(-2)) 从第三个参数开始的参数作为错误日志信息但是FatalOption除外
//	   c.Fatalf(true, "错误", "hh %s", xapi.FatalWithCode(-2), "ss") 这个和上面的执行结果一致 FatalWithCode作为优先级最高{"msg": "错误","code": -2,}
// 以下情况不执行
// c.Fatalf(nil, "错误")  c.Fatalf(false, "错误")

func (c *Ctx) Fatalf(fs ...interface{}) {
	if fs == nil || len(fs) == 0 {
		return
	}

	condition := fs[0] // 第一个参数

	if err, ok := condition.(xerror.XError); ok { // 看是否是error实现
		x := err.XError()
		if x == nil {
			condition = nil
		} else {
			condition = x
		}
	}

	cf := &fatalConfig{
		condition: condition,
		msg:       DefaultErrorMsg,
		code:      CodeError,
	}

	// 特殊处理 多个数据数据的 最后一个是error
	if len(fs) > 1 {
		f := fs[len(fs)-1]
		if f != nil {
			if err, ok := f.(error); ok {
				cf.condition = err
				fatalCondition(c, cf)
				return
			}
			if err, ok := f.(xerror.XError); ok {
				x := err.XError()
				if x == nil {
					cf.condition = nil
				} else {
					cf.condition = x
				}
				fatalCondition(c, cf)
				return
			}
		}
	}
	if condition == nil { // 第一个参数为空 则直接返回
		return
	}

	if len(fs) > 1 {
		msg := fs[1] // 如果第二个参数是string 则认为是错误信息
		if msgStr, ok := msg.(string); ok {
			cf.msg = msgStr
		}
	}
	if len(fs) > 2 {
		l := fs[2] // 如果第三个参数是string 则认为是log信息 从第三个开始是日志参数
		if logStr, ok := l.(string); ok {
			cf.log = fmt.Sprintf(logStr, xarray.Filter(fs[3:], func(i int, item interface{}) bool {
				if item == nil {
					return false
				}
				if _, ok := item.(FatalOption); ok { // FatalOption 不参与日志
					return false
				}
				return true
			})...)
		}
	}

	for _, f := range fs { // 参数中有 option 直接覆盖
		if option, ok := f.(FatalOption); ok {
			option.apply(cf)
		}
	}
	fatalCondition(c, cf)

}

type fatalConfig struct {
	condition interface{}
	msg       string
	log       string
	err       error
	data      interface{}
	code      int
}

type FatalOption interface {
	apply(cfg *fatalConfig)
}
type OptionFunc func(cfg *fatalConfig)

func (f OptionFunc) apply(cfg *fatalConfig) {
	f(cfg)
}

func FatalWithMsg(msg string) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.msg = msg
	})
}

func FatalWithMsgf(msg string, v ...interface{}) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.msg = fmt.Sprintf(msg, v...)
	})
}

func FatalWithLogf(log string, v ...interface{}) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.log = fmt.Sprintf(log, v...)
	})
}

func FatalWithError(err error) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.err = err
	})
}

func FatalWithData(data interface{}) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.data = data
	})
}

func FatalWithCode(code int) FatalOption {
	return OptionFunc(func(cfg *fatalConfig) {
		cfg.code = code
	})
}

func fatalHandle(ctx *Ctx, cf *fatalConfig) {

	//var logMod = xlog.DefaultLogger.Info

	if cf.err != nil {
		cf.log = fmt.Sprintf("msg: %s  |error: %s", cf.msg, cf.err.Error())
		//logMod = xlog.DefaultLogger.Error
	}

	if cf.condition == nil {
		return
	}
	if cf.log != "" {
		//logMod().CallerSkipFrame(1).Fields(xlog.GetFields(ctx)).Msg(cf.log)
	}
	if ctx.test {
		if cf.data != nil {
			indent, err := json.MarshalIndent(&cf.data, "", "  ")
			if err == nil {
				fmt.Printf("test Fatal: %s \n", string(indent))
			} else {
				fmt.Printf("test result:\n%v\n", cf)
			}
		} else {
			fmt.Printf("test result:\n%v\n", cf.msg)
		}
		//panic("test error")
	} else {
		if cf.data != nil {
			panic(cf.data)
		} else {
			panic(cf.msg)
		}
	}
}

func fatalCondition(c *Ctx, cf *fatalConfig) {
	if cf.condition == nil {
		return
	}
	//var logMod = xlog.DefaultLogger.Info
	switch cc := cf.condition.(type) {
	case bool: // true 才执行
		if !cc {
			return
		}
		cf.data = xerror.New(cf.msg).SetCode(cf.code)
	case FatalOption: // 所有执行
	case string: // 所有执行
		cf.data = xerror.New(cc).SetCode(cf.code)
	case xerror.Error:
		if cf.msg != DefaultErrorMsg {
			cc.Msg = cf.msg
		}
		cf.data = cc.SetCode(cf.code)
	case *xerror.Error:
		if cf.msg != DefaultErrorMsg {
			cc.Msg = cf.msg
		}
		cf.data = cc.SetCode(cf.code)
	case Response, *Response:
		cf.data = cc
	default:
		code := xerror.New(cc)
		if cf.msg != DefaultErrorMsg {
			code.Msg = cf.msg
		}
		cf.data = code.SetCode(cf.code)
		//logMod = xlog.DefaultLogger.Error
		if cf.log != "" {
			cf.log = fmt.Sprintf("log: %s | fatal error: %v", cf.log, cc)
		} else {
			cf.log = fmt.Sprintf("fatal error: %v", cc)
		}
	}

	if cf.log != "" {
		if cf.msg != "" && cf.msg != DefaultErrorMsg {
			//logMod().CallerSkipFrame(1).Fields(xlog.GetFields(c)).Msg(fmt.Sprintf("%s | %s", cf.msg, cf.log))
		} else {
			//logMod().CallerSkipFrame(1).Fields(xlog.GetFields(c)).Msg(cf.log)
		}
	}

	if c.test {
		if cf.data != nil {
			indent, err := json.MarshalIndent(&cf.data, "", "  ")
			if err == nil {
				fmt.Printf("test Fatal: %s \n", string(indent))
			} else {
				fmt.Printf("test result:\n%v\n", cf)
			}
		} else {
			fmt.Printf("test result:\n%v\n", cf.msg)
		}
		if c.TraceId != -99 {
			panic("test error")
		}
	} else {
		if cf.data != nil {
			panic(cf.data)
		} else {
			panic(cf.msg)
		}
	}
}
