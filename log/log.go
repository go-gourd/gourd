package log

import (
	"go.uber.org/zap"
)

func Info(msg string) {
	// example
	logger := zap.NewExample()
	logger.Info(msg)

	// Development
	//logger, _ = zap.NewDevelopment()
	//logger.Info("Development")
	//// Production
	//logger, _ = zap.NewProduction()
	//logger.Info("Production")
}
