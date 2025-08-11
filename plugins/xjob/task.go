package xjob

import (
	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/panjf2000/ants/v2"
	"sync"
	"sync/atomic"
	"time"
)

// Task 定义任务结构体
type Task struct {
	ID    string        // 任务唯一标识
	Job   func() error  // 实际执行的函数
	Retry int           // 最大重试次数
	Timeout time.Duration // 任务超时时间
}

// TaskHandler 任务处理器核心结构
type TaskHandler struct {
	pool       *ants.Pool       // ants协程池实例
	taskQueue  chan *Task       // 任务缓冲通道
	wg         sync.WaitGroup   // 等待组控制优雅退出
	panicChan  chan interface{} // panic通知通道
	maxWorkers int              // 最大工作协程数
	closed     int32            // 标记是否已关闭
	tasks      sync.Map         // 正在执行的任务映射
}

// NewTaskHandler 初始化处理器
func NewTaskHandler(maxWorkers int) (*TaskHandler, error) {
	th := &TaskHandler{
		taskQueue:  make(chan *Task, 1000),    // 带缓冲的任务队列
		panicChan:  make(chan interface{}, 1), // panic缓冲通道
		maxWorkers: maxWorkers,
	}

	// 初始化ants池（带预分配和panic处理）
	pool, err := ants.NewPool(
		maxWorkers,
		ants.WithPreAlloc(true),               // 预分配worker提升性能
		ants.WithPanicHandler(th.handlePanic), // 全局panic处理器
	)
	if err != nil {
		return nil, xerror.Newf("failed to create ants pool: %s", err)
	}
	th.pool = pool

	go th.dispatch() // 启动任务分发协程
	return th, nil
}

// Submit 提交任务到处理器
func (th *TaskHandler) Submit(task *Task) error {
	// 检查是否已经关闭
	if atomic.LoadInt32(&th.closed) == 1 {
		return xerror.Newf("task handler is closed")
	}

	// 检查任务ID是否已存在
	if task.ID != "" {
		if _, loaded := th.tasks.LoadOrStore(task.ID, true); loaded {
			return xerror.Newf("task with ID %s already exists", task.ID)
		}
	}

	th.taskQueue <- task // 发送到任务队列
	th.wg.Add(1)         // 增加等待计数
	return nil
}

// dispatch 任务分发器（运行在独立goroutine）
func (th *TaskHandler) dispatch() {
	for task := range th.taskQueue {
		currentTask := task // 闭包捕获当前任务
		// 将 wg.Add 移动到更靠近执行的地方，减少时间差
		err := th.pool.Submit(func() {
			defer th.wg.Done() // 确保任务完成时通知等待组

			// 任务完成后从映射中删除
			defer func() {
				if currentTask.ID != "" {
					th.tasks.Delete(currentTask.ID)
				}
			}()

			// 创建带超时的上下文（如果设置了超时）
			var timeoutChan <-chan time.Time
			if currentTask.Timeout > 0 {
				timer := time.NewTimer(currentTask.Timeout)
				defer timer.Stop()
				timeoutChan = timer.C
			}

			// 指数退避重试逻辑
			for i := 0; i <= currentTask.Retry; i++ {
				// 创建带取消功能的任务执行
				done := make(chan error, 1)
				go func() {
					done <- currentTask.Job()
				}()

				var err error
				if currentTask.Timeout > 0 {
					select {
					case err = <-done:
						// 任务完成
					case <-timeoutChan:
						err = xerror.Newf("task timeout after %v", currentTask.Timeout)
					}
				} else {
					err = <-done
				}

				if err == nil {
					xlog.Tracef(defaultCtx(), "Task %s completed successfully", currentTask.ID)
					return // 成功则退出
				}

				xlog.Warnf(defaultCtx(), "Task %s failed (attempt %d/%d): %v",
					currentTask.ID, i+1, currentTask.Retry+1, err)

				if i < currentTask.Retry {
					// 使用指数退避策略
					time.Sleep(time.Millisecond * time.Duration(1<<uint(i)*100))
				}
			}
		})

		// 如果提交失败，需要减少等待计数
		if err != nil {
			th.wg.Done()
			if currentTask.ID != "" {
				th.tasks.Delete(currentTask.ID)
			}
			xlog.Errorf(defaultCtx(), "Failed to submit task %s: %v", currentTask.ID, err)
		}
	}
}

// CancelTask 取消任务（如果任务支持取消）
func (th *TaskHandler) CancelTask(taskID string) bool {
	_, loaded := th.tasks.LoadAndDelete(taskID)
	return loaded
}

// IsTaskRunning 检查任务是否正在运行
func (th *TaskHandler) IsTaskRunning(taskID string) bool {
	_, ok := th.tasks.Load(taskID)
	return ok
}

// handlePanic 全局panic处理（符合ants.PanicHandler类型）
func (th *TaskHandler) handlePanic(p interface{}) {
	xlog.Errorf(defaultCtx(), "[SYSTEM RECOVER] panic: %v", p)

	// 同时发送到panicChan以供外部监听
	select {
	case th.panicChan <- p:
	default:
		// 防止阻塞
	}
}

// Dispose 等待所有任务完成并释放资源
func (th *TaskHandler) Dispose() error {
	// 使用原子操作检查是否已经关闭
	if !atomic.CompareAndSwapInt32(&th.closed, 0, 1) {
		return xerror.Newf("task handler is already closed")
	}

	xlog.Warnf(defaultCtx(), "task handler is disposing...")

	close(th.taskQueue) // 关闭任务通道
	th.wg.Wait()        // 等待所有任务完成
	th.pool.Release()   // 释放协程池资源
	return nil
}

// GetPanicChan 返回panic通知通道
func (th *TaskHandler) GetPanicChan() <-chan interface{} {
	return th.panicChan
}


