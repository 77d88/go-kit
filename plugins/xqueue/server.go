package xqueue

import (
	"context"
	"github.com/77d88/go-kit/basic/xconfig"
	"github.com/77d88/go-kit/plugins/xlog"
	"github.com/hibiken/asynq"
)

type Server struct {
	*asynq.Server
	consumers map[string]ConsumerFunc
}

var Srv *Server

func NewServer(cfg asynq.Config, consumers map[string]ConsumerFunc) *Server {
	if len(consumers) == 0 {
		return nil
	}
	cfg.Logger = &AsynqLogger{}
	mux := asynq.NewServeMux()
	for key, c := range consumers {
		mux.HandleFunc(key, func(ctx context.Context, task *asynq.Task) error {
			return c(ctx, string(task.Payload()))
		})
	}

	Srv = &Server{
		Server:    asynq.NewServer(getRedisConfig(), cfg),
		consumers: consumers,
	}

	xlog.Infof(nil, "xqueue server start listening... ")
	if err := Srv.Start(mux); err != nil {
		xlog.Fatalf(nil, "could not run server: %v", err)
	}
	return Srv
}

type ConsumerFunc func(ctx context.Context, message string) error

func getRedisConfig() asynq.RedisClientOpt {
	var c RedisConfig
	xconfig.ScanKey("redis", &c)

	return asynq.RedisClientOpt{
		Addr:     c.Addr,
		Password: c.Pass,
		DB:       c.Db,
	}
}

type AsynqLogger struct{}

func (a AsynqLogger) Debug(args ...interface{}) {
	xlog.Debugf(nil, args[0].(string), args[1:]...)
}

func (a AsynqLogger) Info(args ...interface{}) {
	xlog.Infof(nil, args[0].(string), args[1:]...)
}

func (a AsynqLogger) Warn(args ...interface{}) {
	xlog.Warnf(nil, args[0].(string), args[1:]...)
}

func (a AsynqLogger) Error(args ...interface{}) {
	xlog.Errorf(nil, args[0].(string), args[1:]...)
}

func (a AsynqLogger) Fatal(args ...interface{}) {
	xlog.Fatalf(nil, args[0].(string), args[1:]...)
}

func (s *Server) Dispose() error {
	s.Shutdown()
	return nil
}
