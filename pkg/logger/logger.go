package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

type contextKey string

const correlationIDKey contextKey = "correlation_id"

var CorrelationIDKey = correlationIDKey

func NewLogger(level string) *Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapLogger, _ := config.Build()
	return &Logger{zapLogger.Sugar()}
}

func (l *Logger) WithCorrelationID(ctx context.Context) *zap.SugaredLogger {
	if l == nil || l.SugaredLogger == nil {
		return nil
	}
	if ctx == nil {
		return l.SugaredLogger
	}
	if val := ctx.Value(correlationIDKey); val != nil {
		if id, ok := val.(string); ok && id != "" {
			return l.SugaredLogger.With("correlation_id", id)
		}
	}
	return l.SugaredLogger
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}
