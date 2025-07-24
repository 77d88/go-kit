// Copyright 2021 dudaodong@gmail.com. All rights reserved.
// Use of this source code is governed by MIT license

// Package function implements some functions for control the function execution and some is for functional programming.
package xfunc

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// After 创建一个函数，当他被调用n或更多次之后将马上触发fn
func After(n int, fn any) func(args ...any) []reflect.Value {
	// Catch programming error while constructing the closure
	mustBeFunction(fn)

	return func(args ...any) []reflect.Value {
		n--
		if n < 1 {
			return unsafeInvokeFunc(fn, args...)
		}
		return nil
	}
}

// Before 创建一个函数，调用次数不超过n次，之后再调用这个函数，将返回一次最后调用fn的结果
func Before(n int, fn any) func(args ...any) []reflect.Value {
	// Catch programming error while constructing the closure
	mustBeFunction(fn)
	var result []reflect.Value
	return func(args ...any) []reflect.Value {
		if n > 0 {
			result = unsafeInvokeFunc(fn, args...)
		}
		if n <= 0 {
			fn = nil
		}
		n--
		return result
	}
}

// CurryFn is for make curry function
type CurryFn[T any] func(...T) T

// New 创建柯里化函数
// 柯里化（Currying）是函数式编程中的一个概念，它将原本接受多个参数的函数转换为一系列只接受单个参数的函数。
// 换句话说，柯里化把一个多参数函数转换成依次调用的一组单参数函数。
// example:
//
//	func main() {
//	   add := func(a, b int) int {
//	       return a + b
//	   }
//
//	   var addCurry function.CurryFn[int] = func(values ...int) int {
//	       return add(values[0], values[1])
//	   }
//	   add1 := addCurry.New(1)
//
//	   result := add1(2)
//
//	   fmt.Println(result)
//
//	   // Output:
//	   // 3
//	}
func (cf CurryFn[T]) New(val T) func(...T) T {
	return func(vals ...T) T {
		args := append([]T{val}, vals...)
		return cf(args...)
	}
}

// Compose 从右至左组合函数列表fnList，返回组合后的函数
// example:
//
//	func main() {
//		toUpper := func(strs ...string) string {
//			return strings.ToUpper(strs[0])
//		}
//		toLower := func(strs ...string) string {
//			return strings.ToLower(strs[0])
//		}
//		transform := function.Compose(toUpper, toLower)
//
//		result := transform("aBCde")
//
//		fmt.Println(result)
//
//		// Output:
//		// ABCDE
//	}
func Compose[T any](fnList ...func(...T) T) func(...T) T {
	return func(args ...T) T {
		firstFn := fnList[0]
		restFns := fnList[1:]

		if len(fnList) == 1 {
			return firstFn(args...)
		}

		fn := Compose[T](restFns...)
		arg := fn(args...)

		return firstFn(arg)
	}
}

// Delay 延迟delay时间后调用函数
func Delay(delay time.Duration, fn any, args ...any) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		mustBeFunction(fn)

		time.Sleep(delay)
		unsafeInvokeFunc(fn, args...)
	}()
}

// DelaySync 当前线程同步 延迟delay时间后调用函数
func DelaySync(delay time.Duration, fn any, args ...any) {
	mustBeFunction(fn)

	time.Sleep(delay)
	unsafeInvokeFunc(fn, args...)
}

// Debounce 创建一个函数的去抖动版本。该去抖动函数仅在上次调用后的指定延迟时间过去之后才会调用原始函数。支持取消去抖动。
func Debounce(fn func(), delay time.Duration) (debouncedFn func(), cancelFn func()) {
	var (
		timer *time.Timer
		mu    sync.Mutex
	)

	debouncedFn = func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, func() {
			fn()
		})
	}

	cancelFn = func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}
	}

	return debouncedFn, cancelFn
}

// Throttle creates a throttled version of the provided function.
// The returned function guarantees that it will only be invoked at most once per interval.
// Play: https://go.dev/play/p/HpoMov-tJSN
func Throttle(fn func(), interval time.Duration) func() {
	var (
		timer   *time.Timer
		lastRun time.Time
		mu      sync.Mutex
	)

	return func() {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		if now.Sub(lastRun) >= interval {
			fn()
			lastRun = now
			if timer != nil {
				timer.Stop()
				timer = nil
			}
		} else if timer == nil {
			delay := interval - now.Sub(lastRun)

			timer = time.AfterFunc(delay, func() {
				mu.Lock()
				defer mu.Unlock()

				fn()
				lastRun = time.Now()
				timer = nil
			})
		}
	}
}

// Schedule 每次持续时间调用函数，直到关闭返回的 bool chan
func Schedule(duration time.Duration, fn any, args ...any) chan bool {
	// Catch programming error while constructing the closure
	mustBeFunction(fn)

	quit := make(chan bool)

	go func() {
		for {
			unsafeInvokeFunc(fn, args...)

			select {
			case <-time.After(duration):
			case <-quit:
				return
			}
		}
	}()

	return quit
}

// Pipeline 执行函数pipeline.
// example:
//
//	func main() {
//		addOne := func(x int) int {
//			return x + 1
//		}
//		double := func(x int) int {
//			return 2 * x
//		}
//		square := func(x int) int {
//			return x * x
//		}
//
//		fn := function.Pipeline(addOne, double, square)
//
//		result := fn(2)
//
//		fmt.Println(result)
//
//		// Output:
//		// 36
//	}
func Pipeline[T any](funcs ...func(T) T) func(T) T {
	return func(arg T) (result T) {
		result = arg
		for _, fn := range funcs {
			result = fn(result)
		}
		return
	}
}

// AcceptIf returns another function of the same signature as the apply function but also includes a bool value to indicate success or failure.
// A predicate function that takes an argument of type T and returns a bool.
// An apply function that also takes an argument of type T and returns a modified value of the same type.
// Play: https://go.dev/play/p/XlXHHtzCf7d
func AcceptIf[T any](predicate func(T) bool, apply func(T) T) func(T) (T, bool) {
	if predicate == nil {
		panic("programming error: predicate must be not nil")
	}

	if apply == nil {
		panic("programming error: apply must be not nil")
	}

	return func(t T) (T, bool) {
		if !predicate(t) {
			var defaultValue T
			return defaultValue, false
		}
		return apply(t), true
	}
}

func unsafeInvokeFunc(fn any, args ...any) []reflect.Value {
	fv := reflect.ValueOf(fn)
	params := make([]reflect.Value, len(args))
	for i, item := range args {
		params[i] = reflect.ValueOf(item)
	}
	return fv.Call(params)
}

func mustBeFunction(function any) {
	v := reflect.ValueOf(function)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("Invalid function type, value of type %T", function))
	}
}
