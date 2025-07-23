package xjob

import (
	"context"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron/v3"
)

type JobManager struct {
	pool        *ants.PoolWithFunc
	cronManager *CronManager
}

var (
	jobManager *JobManager
)

func getCtx() context.Context {
	return context.WithValue(context.Background(), xlog.CtxLogParam, map[string]interface{}{
		"origin": "xjob",
	})
}

type Logger struct {
}

func (l Logger) Printf(format string, args ...any) {
	xlog.Debugf(getCtx(), format, args...)
}

// Info logs routine messages about cron's operation.
func (l Logger) Info(msg string, keysAndValues ...interface{}) {
	xlog.Infof(getCtx(), msg, keysAndValues...)

}

// Error logs an error condition.
func (l Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	xlog.Errorf(getCtx(), msg, keysAndValues...)
}

func Init() *JobManager {
	// 创建一个具有 100 个 goroutines 的池
	p, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		defer func() {
			if err := recover(); err != nil {
				xlog.Errorf(getCtx(), "panic: %v", err)
			}
		}()
		i.(func())()
	}, ants.WithLogger(Logger{}))
	jobManager = &JobManager{
		pool:        p,
		cronManager: NewCronManager(),
	}
	xlog.Infof(getCtx(), "init xjob success")
	return jobManager
}

func Submit(c context.Context, f func(c context.Context)) error {
	return jobManager.pool.Invoke(func() {
		f(c)
	})
}

func SubmitCron(jobID string, spec string, cmd func()) (cron.EntryID, error) {
	return jobManager.cronManager.AddJob(jobID, spec, cmd)
}

func (d *JobManager) Dispose() error {
	xlog.Warnf(nil, "close xjob")
	d.pool.Release()
	d.cronManager.Stop()
	return nil
}
