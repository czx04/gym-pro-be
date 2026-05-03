package logger

import (
	"fmt"
	"gym-pro-2026-ptit/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface: msg + optional key-value pairs (key, value, key, value, ...).
type Logger interface {
	Debug(msg string, kvs ...interface{})
	Info(msg string, kvs ...interface{})
	Warn(msg string, kvs ...interface{})
	Error(msg string, kvs ...interface{})
	Fatal(msg string, kvs ...interface{})
	With(kvs ...interface{}) Logger
	Sync() error
}

var global Logger = &noopLogger{}

// SetGlobal sets the global logger (bootstrap gọi sau khi New).
func SetGlobal(l Logger) {
	if l != nil {
		global = l
	}
}

// Package-level: gọi trực tiếp logger.Info("msg") hoặc logger.Info("msg", "key", value).
func Debug(msg string, kvs ...interface{}) { global.Debug(msg, kvs...) }
func Info(msg string, kvs ...interface{})  { global.Info(msg, kvs...) }
func Warn(msg string, kvs ...interface{})  { global.Warn(msg, kvs...) }
func Error(msg string, kvs ...interface{}) { global.Error(msg, kvs...) }
func Fatal(msg string, kvs ...interface{}) { global.Fatal(msg, kvs...) }
func Sync() error                          { return global.Sync() }

type noopLogger struct{}

func (n *noopLogger) Debug(string, ...interface{}) {}
func (n *noopLogger) Info(string, ...interface{}) {}
func (n *noopLogger) Warn(string, ...interface{}) {}
func (n *noopLogger) Error(string, ...interface{}) {}
func (n *noopLogger) Fatal(string, ...interface{}) {}
func (n *noopLogger) With(...interface{}) Logger   { return n }
func (n *noopLogger) Sync() error                  { return nil }

type zapLogger struct{ logger *zap.Logger }

func kvsToFields(kvs ...interface{}) []zap.Field {
	if len(kvs) == 0 {
		return nil
	}
	fields := make([]zap.Field, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			key = fmt.Sprintf("key_%d", i/2)
		}
		fields = append(fields, zap.Any(key, kvs[i+1]))
	}
	return fields
}

func (l *zapLogger) Debug(msg string, kvs ...interface{}) { l.logger.Debug(msg, kvsToFields(kvs...)...) }
func (l *zapLogger) Info(msg string, kvs ...interface{})  { l.logger.Info(msg, kvsToFields(kvs...)...) }
func (l *zapLogger) Warn(msg string, kvs ...interface{})  { l.logger.Warn(msg, kvsToFields(kvs...)...) }
func (l *zapLogger) Error(msg string, kvs ...interface{}) { l.logger.Error(msg, kvsToFields(kvs...)...) }
func (l *zapLogger) Fatal(msg string, kvs ...interface{}) { l.logger.Fatal(msg, kvsToFields(kvs...)...) }
func (l *zapLogger) With(kvs ...interface{}) Logger {
	return &zapLogger{logger: l.logger.With(kvsToFields(kvs...)...)}
}
func (l *zapLogger) Sync() error { return l.logger.Sync() }

func New(cfg *config.LoggerConfig) (Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoding := cfg.Encoding
	if encoding == "console" || level == zapcore.DebugLevel {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		encoding = "console"
	} else {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoding = "json"
	}
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       level == zapcore.DebugLevel,
		DisableCaller:     false,
		DisableStacktrace: level != zapcore.ErrorLevel && level != zapcore.FatalLevel,
		Sampling:          nil,
		Encoding:          encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       cfg.OutputPaths,
		ErrorOutputPaths:  cfg.ErrorOutputPaths,
	}
	z, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{logger: z}, nil
}
