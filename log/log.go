package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// example
//logger := zap.NewExample()
// Development
//logger, _ = zap.NewDevelopment()
//// Production
//logger, _ = zap.NewProduction()

func getLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func Info(msg string, fields ...zapcore.Field) {
	logger := getLogger()
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	logger := getLogger()
	logger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger := getLogger()
	logger.Warn(msg, fields...)
}

func Error(err error, fields ...zapcore.Field) {
	logger := getLogger()
	logger.Error(err.Error(), fields...)
}
