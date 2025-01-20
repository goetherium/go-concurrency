package app

import (
	"log/slog"

	"database/internal/entity/parser"
	"database/internal/infrastructure/connhandler"
	"database/internal/infrastructure/engine/hashtable"
	"database/internal/infrastructure/server"
	"database/internal/service/compute"
	"database/internal/service/storage"
	"database/internal/usecase/database"
)

type App struct {
	Logger *slog.Logger
	Server *server.Server
}

func NewApp(logger *slog.Logger) *App {
	p := parser.New()
	computer := compute.New(p)

	engine := hashtable.New(1000)
	store := storage.New(engine)
	db := database.New(computer, store)
	handler := connhandler.New(db)
	server := server.New(logger, handler)

	return &App{
		Logger: logger,
		Server: server,
	}
}
