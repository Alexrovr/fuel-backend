package main

import (
	"github.com/sirupsen/logrus"

	"heat-backend/internal/api"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	if err := api.StartServer(logger); err != nil {
		logger.Fatalf("не удалось запустить сервер: %v", err)
	}
}
