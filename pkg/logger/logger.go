package logger

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
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

type zapLogger struct {
	logger *zap.Logger
}

var _ Interface = (*zapLogger)(nil)

func New(level string, env string) *zapLogger {
	l := new(zapcore.Level)
	switch level {
	case "debug":
		*l = zapcore.DebugLevel
	case "info":
		*l = zapcore.InfoLevel
	case "warn":
	}
	config := new(zap.Config)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	if env == "PROD" || env == "dev" {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(*l),
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"./.log/server.log"},
			ErrorOutputPaths: []string{"./.log/server.log"},
		}

		err := os.MkdirAll("./.log", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(
			"./.log/server.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	} else {
		config = &zap.Config{
			Level:            zap.NewAtomicLevelAt(*l),
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}
	logger := zap.Must(config.Build(zap.AddCallerSkip(1)))

	return &zapLogger{
		logger: logger,
	}
}

func (l *zapLogger) Debug(message string, args ...zap.Field) {
	l.logger.Debug(message, args...)
}

func (l *zapLogger) Info(message string, args ...zap.Field) {
	l.logger.Info(message, args...)
}

func (l *zapLogger) Warn(message string, args ...zap.Field) {
	l.logger.Warn(message, args...)
}

func (l *zapLogger) Error(message error, args ...zap.Field) {
	l.logger.Error(message.Error(), args...)
}

func (l *zapLogger) Fatal(message string, args ...zap.Field) {
	l.logger.Fatal(message, args...)
}
