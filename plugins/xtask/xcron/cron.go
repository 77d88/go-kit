package xcron

/**

CRON 表达式格式
cron 表达式表示一组时间，使用 6 个空格分隔的字段(默认是5个支持分钟级别 开启6个分割需要使用WithSeconds开启秒级支持 一般不开启 通常使用 @every 5s 来进行秒级任务)。

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
示例
0 0 * * * *   每小时执行一次
0 0 0 * * *   每天执行一次


一些预设的cron 表达式。
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

固定 间隔任务
@every <duration>
duration 是可以使用 time.ParseDuration() 解析的格式。
参考 https://pkg.go.dev/time#ParseDuration
示例
@every 5s 代表 五秒执行一次
@every 5m 代表 五分钟执行一次
@every 1h30m 代表 一小时30分执行一次
该间隔不考虑作业运行时。例如，如果作业需要 3 分钟才能运行，并且计划每 5 分钟运行一次，则每次运行之间只有 2 分钟的空闲时间。
*/
import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/77d88/go-kit/basic/xerror"
	"github.com/77d88/go-kit/plugins/x"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/plugins/xtask/xjob"
	"github.com/robfig/cron/v3"
)

var defaultCron *Manager

var defaultCtx = context.WithValue(context.Background(), xlog.CtxLogParam, map[string]interface{}{
	"origin": "xcron",
})

// CronTask 定义定时任务结构体
type CronTask struct {
	ID      string                          // 任务唯一标识
	Spec    string                          // cron 表达式
	Job     func(ctx context.Context) error // 实际执行的函数
	Retry   int                             // 最大重试次数
	Timeout time.Duration                   // 任务超时时间
}

// Manager 定时任务管理器核心结构
type Manager struct {
	cron       *cron.Cron         // cron 实例
	tasks      sync.Map           // 存储任务信息
	entries    sync.Map           // 存储 entry ID 映射
	taskPool   *xjob.Manager      // 任务执行池
	closed     int32              // 标记是否已关闭
	cancelFunc context.CancelFunc // 用于取消所有任务
	own        bool               // 是否开启独立任务池
}

// CronTaskManagerOption 配置选项
type CronTaskManagerOption func(*Manager)

// WithTaskPool 设置任务池
func WithTaskPool(pool *xjob.Manager) CronTaskManagerOption {
	return func(manager *Manager) {
		manager.taskPool = pool
	}
}

// WithSeconds 设置cron 的秒级支持
func WithSeconds() CronTaskManagerOption {
	return func(manager *Manager) {
		cron.WithSeconds()(manager.cron)
	}
}

func NewX() *Manager {
	opts := make([]CronTaskManagerOption, 0, 3)
	i := x.ConfigInt("task.cron.ownWorks") // 独立工作线程数量默认50
	if i == 0 {
		i = 50
	}
	var own = x.ConfigBool("task.cron.own")
	if own { // 是否开启独立处理线程模式 不然使用全局处理线程
		handler, err := xjob.New(i)
		if err != nil {
			xlog.Errorf(nil, "new task handler error: %v", err)
			return nil
		}
		opts = append(opts, WithTaskPool(handler))
	} else {
		get, err := x.Get[*xjob.Manager]()
		if err != nil {
			xlog.Errorf(nil, "get task handler error: %v", err)
			return nil
		}
		opts = append(opts, WithTaskPool(get))
	}
	if x.ConfigBool("task.cron.seconds") { // 是否开启秒级任务
		opts = append(opts, WithSeconds())
	}

	handler, err := New(opts...)
	handler.own = own
	if err != nil {
		xlog.Fatalf(nil, "new cron task manager error: %v", err)
	}
	defaultCron = handler
	return handler
}

// New 创建新的定时任务管理器
func New(opts ...CronTaskManagerOption) (*Manager, error) {
	manager := &Manager{
		cron: cron.New(cron.WithLogger(&Logger{})),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(manager)
	}

	// 如果没有提供任务池，则创建一个默认的
	if manager.taskPool == nil {
		taskHandler, err := xjob.New(10)
		if err != nil {
			return nil, xerror.Newf("failed to create default task handler: %s", err)
		}
		manager.taskPool = taskHandler
	}

	_, cancel := context.WithCancel(context.Background())
	manager.cancelFunc = cancel

	// 启动 cron 调度器
	manager.cron.Start()

	return manager, nil
}

// SubmitCronTask 提交定时任务到管理器
func (cm *Manager) SubmitCronTask(task *CronTask) (cron.EntryID, error) {
	// 检查是否已经关闭
	if atomic.LoadInt32(&cm.closed) == 1 {
		return 0, xerror.Newf("cron task manager is closed")
	}

	// 检查任务ID是否已存在
	if task.ID == "" {
		return 0, xerror.Newf("task ID cannot be empty")
	}

	// 检查是否已存在相同ID的任务
	if _, exists := cm.entries.Load(task.ID); exists {
		return 0, xerror.Newf("cron task with ID %s already exists", task.ID)
	}

	// 包装任务以支持重试和超时
	wrappedJob := cm.wrapJob(task)

	// 添加到 cron 调度器
	entryID, err := cm.cron.AddJob(task.Spec, wrappedJob)
	if err != nil {
		return 0, xerror.Newf("failed to add cron job: %s", err)
	}

	// 存储任务信息
	cm.tasks.Store(entryID, task)
	cm.entries.Store(task.ID, entryID)

	xlog.Tracef(defaultCtx, "Cron task submitted - ID: %s, Spec: %s, EntryID: %d", task.ID, task.Spec, entryID)
	return entryID, nil
}

// wrapJob 包装任务以支持重试和超时
func (cm *Manager) wrapJob(task *CronTask) cron.Job {
	return cron.FuncJob(func() {
		// 为每次执行生成唯一的任务ID
		executionID := fmt.Sprintf("%s-exec-%d", task.ID, time.Now().UnixNano())

		// 提交到任务池执行
		err := cm.taskPool.Submit(&xjob.Task{
			ID:      executionID,
			Job:     task.Job,
			Retry:   task.Retry,
			Timeout: task.Timeout,
		})

		if err != nil {
			xlog.Errorf(defaultCtx, "Failed to submit cron task to pool - TaskID: %s, Error: %v", task.ID, err)
		} else {
			xlog.Tracef(defaultCtx, "Cron task %s submitted to pool with execution ID %s", task.ID, executionID)
		}
	})
}

// RemoveCronTask 移除定时任务
func (cm *Manager) RemoveCronTask(taskID string) error {
	// 检查是否已经关闭
	if atomic.LoadInt32(&cm.closed) == 1 {
		return xerror.Newf("cron task manager is closed")
	}

	// 查找 entry ID
	entryIDValue, ok := cm.entries.Load(taskID)
	if !ok {
		return xerror.Newf("task not found: %s", taskID)
	}

	entryID := entryIDValue.(cron.EntryID)
	cm.cron.Remove(entryID)
	cm.tasks.Delete(entryID)
	cm.entries.Delete(taskID)

	xlog.Infof(defaultCtx, "Cron task removed - TaskID: %s, EntryID: %d", taskID, entryID)
	return nil
}

// GetCronTask 获取定时任务信息
func (cm *Manager) GetCronTask(taskID string) (*CronTask, bool) {
	entryIDValue, ok := cm.entries.Load(taskID)
	if !ok {
		return nil, false
	}

	entryID := entryIDValue.(cron.EntryID)
	taskValue, ok := cm.tasks.Load(entryID)
	if !ok {
		return nil, false
	}

	return taskValue.(*CronTask), true
}

// ListCronTasks 列出所有定时任务
func (cm *Manager) ListCronTasks() []CronTaskInfo {
	var tasks []CronTaskInfo

	cm.entries.Range(func(key, value interface{}) bool {
		taskID := key.(string)
		entryID := value.(cron.EntryID)

		if taskValue, ok := cm.tasks.Load(entryID); ok {
			task := taskValue.(*CronTask)
			entry := cm.cron.Entry(entryID)

			tasks = append(tasks, CronTaskInfo{
				TaskID:  taskID,
				EntryID: entryID,
				Spec:    task.Spec,
				Next:    entry.Next,
				Prev:    entry.Prev,
				Valid:   entry.Valid(),
			})
		}
		return true
	})

	return tasks
}

// CronTaskInfo 定时任务信息
type CronTaskInfo struct {
	TaskID  string       // 任务ID
	EntryID cron.EntryID // Entry ID
	Spec    string       // cron表达式
	Next    time.Time    // 下次执行时间
	Prev    time.Time    // 上次执行时间
	Valid   bool         // 是否有效
}

// Dispose 停止并释放所有资源
func (cm *Manager) Dispose() error {
	// 使用原子操作检查是否已经关闭
	if !atomic.CompareAndSwapInt32(&cm.closed, 0, 1) {
		return xerror.Newf("cron task manager is already closed")
	}

	xlog.Warnf(defaultCtx, "Cron task manager is disposing...")

	// 取消所有任务
	if cm.cancelFunc != nil {
		cm.cancelFunc()
	}

	// 停止 cron 调度器
	cm.cron.Stop()

	// 等待所有任务完成并释放任务池
	if cm.taskPool != nil && cm.own {
		if err := cm.taskPool.Dispose(); err != nil {
			xlog.Errorf(defaultCtx, "Failed to dispose task pool: %v", err)
		}
	}
	xlog.Warnf(defaultCtx, "Cron task manager disposed")

	return nil
}

// GetPanicChan 返回任务池的 panic 通知通道
func (cm *Manager) GetPanicChan() <-chan interface{} {
	return cm.taskPool.GetPanicChan()
}

// IsRunning 检查管理器是否正在运行
func (cm *Manager) IsRunning() bool {
	return atomic.LoadInt32(&cm.closed) == 0
}

// IsCronTaskScheduled 检查定时任务是否已调度
func (cm *Manager) IsCronTaskScheduled(taskID string) bool {
	_, ok := cm.entries.Load(taskID)
	return ok
}

type Logger struct {
}

func (l Logger) Printf(format string, args ...any) {
	xlog.Tracef(defaultCtx, format, args...)
}

// Info logs routine messages about cron's operation.
func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	xlog.Tracef(defaultCtx, msg, keysAndValues...)
}

// Error logs an error condition.
func (l Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	xlog.Errorf(defaultCtx, msg, keysAndValues...)
}

func Submit(task *CronTask) error {
	_, err := defaultCron.SubmitCronTask(task)
	return err
}
