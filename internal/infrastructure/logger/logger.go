package logger

import (
	"gym-pro-2026-ptit/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface defines logging methods
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

// zapLogger wraps zap.Logger to implement Logger interface
type zapLogger struct {
	logger *zap.Logger
}

// New creates a new logger instance based on config
func New(cfg *config.LoggerConfig) (Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create encoder based on encoding type
	// encoder := zapcore.NewJSONEncoder(encoderConfig)
	// if cfg.Encoding == "console" {
	// 	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	// }

	// Create zap config
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       level == zapcore.DebugLevel,
		DisableCaller:     false,
		DisableStacktrace: level != zapcore.ErrorLevel && level != zapcore.FatalLevel,
		Sampling:          nil,
		Encoding:          cfg.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       cfg.OutputPaths,
		ErrorOutputPaths:  cfg.ErrorOutputPaths,
	}

	// Build logger
	logger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
