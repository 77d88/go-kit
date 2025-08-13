package pgmq

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/77d88/go-kit/plugins/xdatabase/xdb"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/77d88/go-kit/plugins/xtask/xjob"
)

// Queue states
const (
	QueueStatePending = iota // 待处理
	QueueStateRunning        // 已处理
	QueueStateAck            // 已处理
	QueueStateNACK           // 处理失败
)

type MsgType int

const DefaultMsgType = MsgType(0)

// Handler 处理器结构
type Handler func(msg *Queue) (bool, error)

// XQueue 队列主服务
type XQueue struct {
	db             *xdb.DB
	handlers       map[MsgType][]Handler
	runExecutor    *time.Ticker
	executorThread *xjob.TaskHandler
	stopChan       chan struct{}
	mu             sync.RWMutex
}

type Config struct {
	DBName   string // 数据库名称 没有就默认
	MaxWorks int    // 最大处理任务队列
}

// NewMq 创建新的队列服务实例
func NewMq(cfg *Config) *XQueue {

	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.DBName == "" {
		cfg.DBName = xdb.DefaultDbLinkStr
	}
	if cfg.MaxWorks <= 0 {
		cfg.MaxWorks = 50
	}

	handler, err := xjob.NewTaskHandler(cfg.MaxWorks) // 创建任务处理器
	if err != nil {
		panic(err)
	}
	db, err := xdb.GetDB(cfg.DBName)
	if err != nil {
		panic(err)
	}
	return &XQueue{
		db:             db,
		handlers:       make(map[MsgType][]Handler),
		stopChan:       make(chan struct{}),
		executorThread: handler, // 缓冲通道模拟线程池
	}
}

// RegisterHandler 注册队列处理器
func (xq *XQueue) RegisterHandler(msgType MsgType, handler Handler) {
	xq.mu.Lock()
	defer xq.mu.Unlock()

	handlers, e := xq.handlers[msgType]
	if !e {
		handlers = make([]Handler, 0, 10)
	}
	handlers = append(handlers, handler)
	xq.handlers[msgType] = handlers
}

// Start 启动队列服务
func (xq *XQueue) Start() error {
	xq.mu.RLock()
	hasHandlers := len(xq.handlers) > 0
	xq.mu.RUnlock()

	if !hasHandlers {
		return nil
	}

	// 启动检查任务
	go xq.startCheckTask()

	// 启动消费任务
	go xq.startConsumptionTask()

	xlog.Infof(context.Background(), "Start queue successful handlers: %s", func() string {
		keys := make([]string, 0, len(xq.handlers))
		for k, v := range xq.handlers {
			keys = append(keys, fmt.Sprintf("type:[%d]:handler(%d)", k, len(v)))
		}
		return strings.Join(keys, "\n")
	}())

	return nil
}

// startCheckTask 启动检查任务
func (xq *XQueue) startCheckTask() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			xq.checkTimeoutQueues()
		case <-xq.stopChan:
			return
		}
	}
}

// checkTimeoutQueues 检查超时队列
func (xq *XQueue) checkTimeoutQueues() {
	defer func() {
		if err := recover(); err != nil {
			xlog.Errorf(context.Background(), "检查超时队列异常: %v", recover())
		}
	}()
	// 查询超时或ACK状态的队列  执行时间超过10分钟认为是 执行超时 如果还有次数进入待执行延迟5秒 如果没有则进入失败状态
	result :=
		xq.db.Exec(`WITH queue AS (
    SELECT * FROM s_queue 
    WHERE (state = 1 AND now() > read_time + interval '10 minute')
    FOR UPDATE SKIP LOCKED
)
UPDATE s_queue SET 
    state = CASE WHEN num < retry THEN 0 ELSE 3 END,
    error_info = 'timeout',
    delivery_time = CURRENT_TIMESTAMP + interval '5 second'
WHERE id IN (SELECT id FROM queue)
	`)
	if result.Error != nil {
		log.Printf("Failed to query timeout queues: %v", result.Error)
		return
	}

}

// startConsumptionTask 启动消费任务
func (xq *XQueue) startConsumptionTask() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			xq.consumeQueues()
		case <-xq.stopChan:
			return
		}
	}
}

// consumeQueues 消费队列消息
func (xq *XQueue) consumeQueues() {

	// 查询待处理的队列消息
	var queues []*Queue
	result := xq.db.WithCtx(context.Background()).Raw(`
		WITH queue AS (
			SELECT * FROM s_queue 
			WHERE state = 0 AND delivery_time < now() AND num < retry 
			LIMIT 5 
			FOR UPDATE SKIP LOCKED
		)
		UPDATE s_queue 
		SET state = 1, read_time = now(), num = num + 1 
		WHERE id IN (SELECT id FROM queue)
		RETURNING s_queue.*
	`).Find(&queues)
	if result.Error != nil {
		return
	}

	// 异步处理队列消息
	for _, queue := range queues {
		q := queue
		err := xq.executorThread.Submit(&xjob.Task{
			ID: strconv.FormatInt(q.ID, 10),
			Job: func(ctx context.Context) error {
				return xq.processQueue(ctx, q)
			},
			Retry:   0,
			Timeout: time.Minute * 5,
			Ctx: context.WithValue(context.Background(), xlog.CtxLogParam, map[string]interface{}{
				"QID":   q.ID,
				"QTYPE": q.Type,
			}),
		})
		if err != nil { // 这种失败的重试
			xq.updateQueueState(q.ID, QueueStatePending, err.Error())
		}
	}
}

// processQueue 处理单个队列消息
func (xq *XQueue) processQueue(ctx context.Context, queue *Queue) error {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			errorInfo := string(buf[:n])
			xq.updateQueueState(queue.ID, QueueStateNACK, errorInfo)
		}
	}()

	handlerList, exists := xq.handlers[MsgType(queue.Type)]

	if !exists || len(handlerList) == 0 {
		xlog.Warnf(ctx, "Queue %d -> %d has no handler", queue.ID, queue.Type)
		if queue.Num < queue.Retry {
			xq.updateQueueState(queue.ID, QueueStatePending, "")
		} else {
			xq.updateQueueState(queue.ID, QueueStateNACK, "")
		}
		return nil
	}

	// 执行所有处理器
	for _, handler := range handlerList {
		ack, err := handler(queue)
		if err != nil {
			xlog.Warnf(ctx, "Handler execution failed for queue %d: %v", queue.ID, err)
		}
		if !ack {
			xlog.Warnf(ctx, "Handler execution nack for queue %d: %v", queue.ID, err)
		}
		if !ack || err != nil {
			if queue.Num < queue.Retry {
				xq.updateQueueState(queue.ID, QueueStatePending, err.Error())
			} else {
				xq.updateQueueState(queue.ID, QueueStateNACK, err.Error())
			}
			break // 不再处理下面的
		}
	}
	// 所有处理器执行成功
	xq.updateQueueState(queue.ID, QueueStateAck, "")
	return nil
}

// updateQueueState 更新队列状态
func (xq *XQueue) updateQueueState(queueID int64, state int, errorInfo string) {
	if result := xq.db.Exec("UPDATE s_queue SET state = ?, error_info = ? WHERE id = ?", state, errorInfo, queueID); result.Error != nil {
		xlog.Errorf(context.Background(), "更新队列状态异常: %v", result.Error)
		return
	}
}

// Stop 停止队列服务
func (xq *XQueue) Stop() error {
	close(xq.stopChan)

	return nil
}
