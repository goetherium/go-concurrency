package app

import (
	"log/slog"

	"database/internal/config"
	"database/internal/entity/parser"
	"database/internal/infrastructure/connhandler"
	"database/internal/infrastructure/engine/hashtable"
	"database/internal/infrastructure/semaphore"
	"database/internal/infrastructure/server"
	"database/internal/service/compute"
	"database/internal/service/storage"
	"database/internal/usecase/database"
)

type App struct {
	Logger *slog.Logger
	Server *server.Server
}

func NewApp(logger *slog.Logger, cfg *config.App) *App {
	p := parser.New()
	computer := compute.New(p)

	engine := hashtable.New(1000)
	store := storage.New(engine)
	db := database.New(computer, store)

	handler := connhandler.New(db, cfg.Connection.IdleTimeout, cfg.Connection.MaxMessageSize)
	sema := semaphore.New(cfg.Connection.MaxConnections, cfg.Connection.ConnectTimeout)
	s := server.New(logger, sema, handler)

	return &App{
		Logger: logger,
		Server: s,
	}
}
