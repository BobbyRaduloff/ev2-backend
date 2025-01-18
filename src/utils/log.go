package utils

import (
	"log"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()

	if err != nil {
		log.Fatalf("can't initialize zap logger.\n")
	}

	Logger = logger
}

func SyncLogger(timeout time.Duration) {
	for {
		time.Sleep(timeout)
		Logger.Sync()
	}
}
