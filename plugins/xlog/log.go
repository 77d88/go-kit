package xlog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const CtxLogParam = "_xe/CTX_LOG_PARAM"

type Func = func(f LoggerFields, format string, v ...interface{})
type Logger struct {
	*zap.Logger
	IsTrace bool
}
type LoggerFields interface {
	Fields() Fields
}

var DefaultLogger Logger
var callMap = make(map[string]string)

func init() {
	WithDebugger()
}

var EmptyFields []zap.Field

type Fields map[string]interface{}

func (f Fields) Fields() Fields {
	return f
}

func WithDebugger() {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(fmt.Sprintf("[%s]", level.CapitalString()))
	}
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.TimeOnly))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel, // 测试环境开启Debug级别
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.Development(),                 // 开发模式标记
		zap.AddStacktrace(zap.PanicLevel), // 降低堆栈跟踪阈值
		zap.AddCallerSkip(1),
	)
	defer logger.Sync()
	DefaultLogger = Logger{Logger: logger, IsTrace: false}
}

func WithRelease() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	//encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // 启用彩色日志级别

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), // 使用易读的控制台格式
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel),
		zap.AddCallerSkip(1),
	)
	defer logger.Sync()
	DefaultLogger = Logger{Logger: logger, IsTrace: false}
}

func GetFields(c context.Context) []zap.Field {
	if c == nil {
		return EmptyFields
	}
	f := c.Value(CtxLogParam)
	if f == nil {
		return EmptyFields
	}
	return []zap.Field{zap.Reflect("ext", f)}
}

func (l Logger) Trace(format string, fields ...zap.Field) {
	if !l.IsTrace {
		return
	}
	fields = append(fields, zap.Field{Key: "trace", Type: zapcore.StringType, String: "t"})
	l.Debug(format, fields...)
}

func Tracef(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Trace(fmt.Sprintf(format, v...), GetFields(f)...)
}
func Debugf(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Debug(fmt.Sprintf(format, v...), GetFields(f)...)
}

func Infof(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Info(fmt.Sprintf(format, v...), GetFields(f)...)
}

func Warnf(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Warn(fmt.Sprintf(format, v...), GetFields(f)...)
}

func Errorf(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Error(fmt.Sprintf(format, v...), GetFields(f)...)
}

func Panicf(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Panic(fmt.Sprintf(format, v...), GetFields(f)...)
}

func Fatalf(f context.Context, format string, v ...interface{}) {
	DefaultLogger.Fatal(fmt.Sprintf(format, v...), GetFields(f)...)
}
