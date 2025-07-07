package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Interface interface {
	Debug(message string, args ...zap.Field)
	Info(message string, args ...zap.Field)
	Warn(message string, args ...zap.Field)
	Error(message error, args ...zap.Field)
	Fatal(message string, args ...zap.Field)
}
type Logger struct {
	logger *zap.Logger
}

func New(level string) *Logger {
	l := new(zapcore.Level)
	switch level {
	case "debug":
		*l = zapcore.DebugLevel
	case "info":
		*l = zapcore.InfoLevel
	case "warn":
	}
	env := os.Getenv("ENV")
	config := new(zap.Config)
	if env == "PROD" {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(*l),
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"./.log/server.log"},
			ErrorOutputPaths: []string{"./.log/server.log"},
		}
	} else {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(*l),
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}
	logger := zap.Must(config.Build())

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(message string, args ...zap.Field) {
	l.logger.Debug(message, args...)
}

func (l *Logger) Info(message string, args ...zap.Field) {
	l.logger.Info(message, args...)
}

func (l *Logger) Warn(message string, args ...zap.Field) {
	l.logger.Warn(message, args...)
}

func (l *Logger) Error(message error, args ...zap.Field) {
	l.logger.Error(message.Error(), args...)
}

func (l *Logger) Fatal(message string, args ...zap.Field) {
	l.logger.Fatal(message, args...)
}
