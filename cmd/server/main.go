package main

import (
	"log/slog"
	"os"
	"os/signal"

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

	go func() {
		if err := myApp.Server.ListenAndServe(cfg.Connection.Address); err != nil {
			panic(err)
		}
	}()

	logger.Info("service started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	sig := <-stop

	logger.Info("stopping app...", slog.String("signal", sig.String()))
	logger.Info("service stopped")
}
