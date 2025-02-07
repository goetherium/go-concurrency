package database

import (
	"log/slog"

	"database/internal/entity/command"
	"database/internal/entity/database"
	"database/internal/entity/storage"
)

type computeLayer interface {
	ParseCmd(logger *slog.Logger, rowCmd command.RawCmd) (command.Cmd, error)
}

type storageLayer interface {
	ExecCmd(logger *slog.Logger, cmd command.Cmd) (storage.ExecResult, error)
}

type Database struct {
	compute computeLayer
	storage storageLayer
}

func New(compute computeLayer, storage storageLayer) *Database {
	return &Database{
		compute: compute,
		storage: storage,
	}
}

func (d Database) HandleCmd(logger *slog.Logger, rowCmd command.RawCmd) (database.HandleCmdResult, error) {
	query, err := d.compute.ParseCmd(logger, rowCmd)
	if err != nil {
		return database.HandleCmdResult{}, err
	}

	return d.execCmd(logger, query)
}

func (d Database) execCmd(logger *slog.Logger, cmd command.Cmd) (database.HandleCmdResult, error) {
	execRes, err := d.storage.ExecCmd(logger, cmd)
	if err != nil {
		return database.HandleCmdResult{}, err
	}

	if execRes.Value != "" {
		return database.HandleCmdResult{Result: execRes.Value.String()}, nil
	}

	return database.HandleCmdResult{Result: execRes.Result}, nil
}
