package xapi

//
//import (
//	"context"
//	"github.com/77d88/go-kit/v2/basic/xreflect"
//	"github.com/77d88/go-kit/v2/plugins/xlog"
//	"go.uber.org/dig"
//	"reflect"
//	"sync"
//)
//
//type ApiRes struct {
//	Container *dig.Container
//	disposes  []LifeCycleDispose
//	afterRes  []interface{}
//
//	resLock sync.Mutex
//	resWait sync.WaitGroup
//}
//
//func newRes() *ApiRes {
//	Container := dig.New()
//	return &ApiRes{Container: Container}
//}
//
//// Add 添加资源 自动填充Scanner 如果返回的是实现了LifeCycleDispose接口则自动添加到dispose中 再api结束时自动释放
//func (x *ApiRes) Add(res interface{}) *ApiRes {
//	x.resWait.Add(1)
//	x.add(res)
//	return x
//}
//
//func (x *ApiRes) AddAfter(res ...interface{}) *ApiRes {
//	x.afterRes = append(x.afterRes, res...)
//	return x
//}
//
//func (x *ApiRes) add(dispose interface{}) {
//
//	go func(r interface{}, x *Engine) {
//		defer func() {
//			if err := recover(); err != nil {
//				xlog.Errorf(nil, "res init panic: %v", err)
//			}
//			x.resWait.Done()
//		}()
//		if r == nil {
//			return
//		}
//
//		switch v := r.(type) {
//		case func(*Engine):
//			v(x)
//			return
//		case func(*Engine) interface{}:
//			addStop(x, v(x))
//		case func(*Engine) (interface{}, error):
//			result, err := v(x)
//			if err != nil {
//				xlog.Errorf(nil, "res init error: %v", err)
//			} else {
//				if result != nil {
//					addStop(x, result)
//				}
//			}
//		case LifeCycleDispose:
//			x.disposes = append(x.disposes, v)
//			return
//		case func():
//			v()
//			return
//		case func() interface{}:
//			addStop(x, v())
//		case func() (interface{}, error):
//			result, err := v()
//			if err != nil {
//				xlog.Errorf(nil, "res init error: %v", err)
//			} else {
//				if result != nil {
//					addStop(x, result)
//				}
//			}
//		default:
//			val := reflect.ValueOf(v)
//			if val.Kind() != reflect.Func {
//				addStop(x, v)
//			} else {
//				// 仅支持无参函数
//				typ := val.Type()
//				numIn := typ.NumIn()
//				if numIn != 0 {
//					xlog.Errorf(nil, "res init error: %s func must be no param", val.String())
//					return
//				}
//
//				// 调用函数并处理返回值
//				results := val.Call(make([]reflect.Value, 0))
//				for _, result := range results {
//					if result.Interface() != nil {
//						addStop(x, result.Interface())
//					}
//				}
//			}
//		}
//	}(dispose, x)
//}
//
//func addStop(api *Engine, ress ...interface{}) {
//	if len(ress) == 0 {
//		return
//	}
//	for _, res := range ress {
//		if xreflect.IsNil(res) {
//			continue
//		}
//		if xreflect.ImplementsInterface(res, (*LifeCycleDispose)(nil)) {
//			api.addStopAfter(res.(LifeCycleDispose))
//		}
//	}
//}
//
//func (e *Engine) closeDisposes(ctx context.Ctx) {
//	dwg := &sync.WaitGroup{}
//	done := make(chan struct{}) // 用于通知WaitGroup完成的通道
//	dwg.Add(len(e.disposes))
//
//	go func() {
//		defer close(done) // 确保通道最终关闭
//		for i, dispose := range e.disposes {
//			go func(idx int, d LifeCycleDispose) {
//				defer func() {
//					dwg.Done()
//					if err := recover(); err != nil {
//						xlog.Errorf(nil, "dispose %d panic: %v", idx, err)
//					}
//				}()
//				if !xreflect.IsNil(d) {
//					err := dispose.Dispose()
//					if err != nil {
//						xlog.Errorf(nil, "dispose %d error %v", i, err)
//					}
//				}
//			}(i, dispose)
//		}
//	}()
//	dwg.Wait() // 等待所有goroutine完成
//
//	select {
//	case <-done: // 所有dispose操作完成
//		return
//	case <-ctx.Done(): // 上下文超时或取消
//		panic("App Shutdown timeout")
//	}
//}
