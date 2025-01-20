package main

import (
	"log/slog"

	"database/internal/infrastructure/client"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	zapL := zap.Must(zap.NewDevelopment())

	defer func() {
		_ = zapL.Sync()
	}()

	logger := slog.New(zapslog.NewHandler(zapL.Core()))

	client.Run(logger, ":8080")
}
