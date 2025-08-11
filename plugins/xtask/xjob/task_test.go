package xjob

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 性能测试配置
const (
	benchmarkTaskCount = 100000
	benchmarkWorkers   = 1000
)

// BenchmarkTaskHandlerHighLoad 高负载情况下的性能测试
func BenchmarkTaskHandlerHighLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		handler, err := NewTaskHandler(benchmarkWorkers)
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()
		for j := 0; j < benchmarkTaskCount; j++ {
			taskID := j
			_ = handler.Submit(&Task{
				ID: strconv.Itoa(taskID),
				Job: func() error {
					// 模拟实际工作负载
					time.Sleep(time.Nanosecond * 100)
					return nil
				},
				Retry: 1,
			})
		}

		err = handler.Dispose()
		if err != nil {
			panic(err)
		}
		b.StopTimer()
	}
}

// BenchmarkConcurrency 测试不同并发级别的性能
func BenchmarkConcurrency(b *testing.B) {
	concurrentLevels := []int{10, 50, 100, 200, 500}

	for _, workers := range concurrentLevels {
		b.Run(fmt.Sprintf("Workers-%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				handler, err := NewTaskHandler(workers)
				if err != nil {
					b.Fatal(err)
				}

				b.StartTimer()
				var wg sync.WaitGroup
				// 分批提交任务以避免阻塞
				for batch := 0; batch < benchmarkTaskCount; batch += 1000 {
					wg.Add(1)
					go func(start, end int) {
						defer wg.Done()
						for j := start; j < end && j < benchmarkTaskCount; j++ {
							taskID := j
							_ = handler.Submit(&Task{
								ID: strconv.Itoa(taskID),
								Job: func() error {
									time.Sleep(time.Microsecond)
									return nil
								},
								Retry: 0,
							})
						}
					}(batch, batch+1000)
				}
				wg.Wait()

				err = handler.Dispose()
				if err != nil {
					panic(err)
				}
				b.StopTimer()
			}
		})
	}
}

func TestTask(t *testing.T) {
	// 示例使用
	handler, err := NewTaskHandler(10) // 初始化10个worker的处理器
	if err != nil {
		panic(fmt.Errorf("failed to create task handler: %w", err))
	}

	// 提交会panic的任务（自动恢复）
	err = handler.Submit(&Task{
		ID: "1",
		Job: func() error {
			fmt.Println("执行危险任务...")
			panic("模拟意外崩溃")
		},
		Retry: 2,
	})
	if err != nil {
		fmt.Printf("提交任务失败: %v\n", err)
	}

	// 提交正常任务
	err = handler.Submit(&Task{
		ID: "2",
		Job: func() error {
			fmt.Println("执行安全任务")
			return nil
		},
	})
	if err != nil {
		fmt.Printf("提交任务失败: %v\n", err)
	}

	// 监听panic事件
	go func() {
		for p := range handler.GetPanicChan() {
			fmt.Printf("收到panic通知: %v\n", p)
		}
	}()

	err = handler.Dispose() // 阻塞直到所有任务完成
	if err != nil {
		fmt.Printf("等待任务完成时出错: %v\n", err)
	}
	err = handler.Submit(&Task{
		ID: "3",
		Job: func() error {
			fmt.Println("结束后执行任务")
			return nil
		},
	})
	if err != nil {
		fmt.Printf("提交任务失败: %v\n", err)
	}
}
