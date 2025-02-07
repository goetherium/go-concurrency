package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"database/internal/app"
	"database/internal/config"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	if len(os.Args) < 2 {
		panic("config file not specified in command line")
	}

	cfg := config.Setup(os.Args[1])

	zapLogger := app.NewLogger(cfg.Logger.Level)
	defer func() {
		_ = zapLogger.Sync()
	}()

	logger := slog.New(zapslog.NewHandler(zapLogger.Core()))

	myApp := app.NewApp(logger, &cfg)

	stopCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := myApp.Start(stopCtx); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	go func() {
		if err := myApp.Server.ListenAndServe(cfg.Connection.Address); err != nil {
			panic(err)
		}
	}()

	logger.Info("service started")

	<-stopCtx.Done()

	logger.Info("stopping app...")

	if err := myApp.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown app: %v", err)
	}

	logger.Info("service stopped")
}
