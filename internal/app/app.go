package app

import (
	"context"
	"log/slog"

	"database/internal/config"
	"database/internal/entity/parser"
	"database/internal/infrastructure/connhandler"
	"database/internal/infrastructure/engine/hashtable"
	"database/internal/infrastructure/semaphore"
	"database/internal/infrastructure/server"
	"database/internal/infrastructure/wal"
	"database/internal/service/compute"
	"database/internal/service/storage"
	"database/internal/usecase/database"
)

type App struct {
	Logger  *slog.Logger
	Server  *server.Server
	storage *storage.Storage
}

func NewApp(logger *slog.Logger, cfg *config.App) *App {
	p := parser.New()
	computer := compute.New(p)

	w := wal.NewWal(logger, &cfg.Wal, p)
	engine := hashtable.New(1000)
	store := storage.New(logger, w, engine)
	db := database.New(computer, store)

	handler := connhandler.New(db, cfg.Connection.IdleTimeout, cfg.Connection.MaxMessageSize)
	sema := semaphore.New(cfg.Connection.MaxConnections, cfg.Connection.ConnectTimeout)
	s := server.New(logger, sema, handler)

	return &App{
		Logger:  logger,
		Server:  s,
		storage: store,
	}
}

func (a *App) Start(ctx context.Context) error {
	return a.storage.Open(ctx)
}

func (a *App) Shutdown() error {
	return a.storage.Shutdown()
}
