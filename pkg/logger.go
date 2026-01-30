package pkg

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const TraceIdKey = "traceId"

// Logger 日志结构体
type Logger struct {
	sugaredLogger *zap.SugaredLogger
}

var globalLogger *Logger

func InitLogger(config LogConf) {
	//初始化日志逻辑
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		Compress:   config.Compress,
	}

	writer := zapcore.AddSync(lumberJackLogger)

	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	//同时输出到文件和控制台
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, writer, level),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	globalLogger = &Logger{
		sugaredLogger: logger.Sugar(),
	}
}

func GetLogger() *Logger {
	return globalLogger
}

// WithContext添加traceId到日志
func (l *Logger) WithContext(ctx context.Context) *Logger {
	traceId := ctx.Value(TraceIdKey)
	if traceId != nil {
		newLogger := l.sugaredLogger.With(TraceIdKey, traceId)
		return &Logger{sugaredLogger: newLogger}
	}
	return l
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

// DPanic logs a message at DPanicLevel.
func (l *Logger) DPanic(args ...interface{}) {
	l.sugaredLogger.DPanic(args...)
}

// Panic logs a message at PanicLevel.
func (l *Logger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

// Fatal logs a message at FatalLevel.
func (l *Logger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message at DebugLevel.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugaredLogger.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message at InfoLevel.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message at WarnLevel.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message at ErrorLevel.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.sugaredLogger.Sync()
}
