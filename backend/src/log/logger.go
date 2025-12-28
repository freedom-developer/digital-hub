package logger

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

const (
	LogKeyTag = "tag"
)

type LogConfig struct {
	LogDir     string
	LogLevel   int // 1. debug, 2. info, 3. warn, 4. error, 5. fatal
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		LogDir:     "./logs",
		LogLevel:   int(zapcore.DebugLevel), // 1
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

func parseLogLevel(level int) zapcore.Level {
	switch level {
	case 1:
		return zapcore.DebugLevel
	case 2:
		return zapcore.InfoLevel
	case 3:
		return zapcore.WarnLevel
	case 4:
		return zapcore.ErrorLevel
	case 5:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func InitLogger(config *LogConfig) *zap.SugaredLogger {
	if config == nil {
		config = DefaultLogConfig()
	}

	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		panic("创建日志目录失败: " + err.Error())
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		CallerKey:  "location",
		MessageKey: "msg",
		// 完全不设置 StacktraceKey
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	businessWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "business.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "error.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	})

	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	businessLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= parseLogLevel(config.LogLevel) && lvl < zapcore.ErrorLevel
	})
	businessCore := zapcore.NewCore(
		jsonEncoder,
		zapcore.NewMultiWriteSyncer(businessWriter, stdoutSyncer),
		businessLevel,
	)

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	errorCore := zapcore.NewCore(
		jsonEncoder,
		zapcore.NewMultiWriteSyncer(errorWriter, stderrSyncer),
		errorLevel,
	)

	core := zapcore.NewTee(businessCore, errorCore)

	// 只添加 AddCaller，什么都不加
	logger = zap.New(core, zap.AddCaller()).Sugar()

	return logger
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.TrimmedPath())
}

func getTagFromContext(ctx *context.Context) string {
	tagVal := (*ctx).Value(LogKeyTag)
	if tagVal == nil {
		return ""
	}
	return tagVal.(string)
}

func SetTagInContext(ctx *context.Context, tag string) *context.Context {
	newCtx := context.WithValue(*ctx, LogKeyTag, tag)
	return &newCtx
}

func ZDebug(ctx *context.Context, msg string, keysAndValues ...interface{}) {
	tagStr := getTagFromContext(ctx)
	if tagStr != "" {
		keysAndValues = append([]interface{}{LogKeyTag, tagStr}, keysAndValues...)
	}
	logger.Debugw(msg, keysAndValues...)
}
func ZInfo(ctx *context.Context, msg string, keysAndValues ...interface{}) {
	tagStr := getTagFromContext(ctx)
	if tagStr != "" {
		keysAndValues = append([]interface{}{LogKeyTag, tagStr}, keysAndValues...)
	}
	logger.Infow(msg, keysAndValues...)
}
func ZWarn(ctx *context.Context, msg string, keysAndValues ...interface{}) {
	tagStr := getTagFromContext(ctx)
	if tagStr != "" {
		keysAndValues = append([]interface{}{LogKeyTag, tagStr}, keysAndValues...)
	}
	logger.Warnw(msg, keysAndValues...)
}
func ZError(ctx *context.Context, msg string, err error, keysAndValues ...interface{}) {
	tagStr := getTagFromContext(ctx)
	if tagStr != "" {
		keysAndValues = append([]interface{}{LogKeyTag, tagStr}, keysAndValues...)
	}
	if err != nil {
		keysAndValues = append([]interface{}{"error", err}, keysAndValues...)
	}
	logger.Errorw(msg, keysAndValues...)
}
func ZFatal(ctx *context.Context, msg string, err error, keysAndValues ...interface{}) {
	tagStr := getTagFromContext(ctx)
	if tagStr != "" {
		keysAndValues = append([]interface{}{LogKeyTag, tagStr}, keysAndValues...)
	}
	if err != nil {
		keysAndValues = append([]interface{}{"error", err}, keysAndValues...)
	}
	logger.Fatalw(msg, keysAndValues...)
}
