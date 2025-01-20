package main

import (
	"log/slog"
	"os"
	"os/signal"

	"database/internal/app"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	zapL := zap.Must(zap.NewDevelopment())

	defer func() {
		_ = zapL.Sync()
	}()

	logger := slog.New(zapslog.NewHandler(zapL.Core()))

	myApp := app.NewApp(logger)

	go func() {
		if err := myApp.Server.ListenAndServe(":8080"); err != nil {
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
