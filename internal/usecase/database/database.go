package database

import (
	"log/slog"

	"database/internal/entity/command"
	"database/internal/entity/database"
	"database/internal/entity/logmodel"
	"database/internal/entity/storage"
)

type computeLayer interface {
	ParseCmd(logger *slog.Logger, rowCmd string) (command.Query, error)
}

type storageLayer interface {
	Set(logger *slog.Logger, key command.Key, value command.Value) (storage.ExecResult, error)
	Get(logger *slog.Logger, key command.Key) (storage.ExecResult, error)
	Del(logger *slog.Logger, key command.Key) (storage.ExecResult, error)
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

func (d Database) HandleCmd(logger *slog.Logger, rowCmd string) (database.HandleCmdResult, error) {
	query, err := d.compute.ParseCmd(logger, rowCmd)
	if err != nil {
		return database.HandleCmdResult{}, err
	}

	return d.execCmd(logger, query)
}

func (d Database) execCmd(logger *slog.Logger, query command.Query) (database.HandleCmdResult, error) {
	child := logger.With(slog.String(
		logmodel.FieldAction, "ExecCmd"),
		"cmdID", query.CmdID,
		"args", query.Args,
	)

	var (
		execRes storage.ExecResult
		err     error
	)

	switch query.CmdID {
	case command.SetCommandID:
		execRes, err = d.storage.Set(logger, command.Key(query.Args[0]), command.Value(query.Args[1]))
	case command.GetCommandID:
		execRes, err = d.storage.Get(logger, command.Key(query.Args[0]))
	case command.DelCommandID:
		execRes, err = d.storage.Del(logger, command.Key(query.Args[0]))
	default:
		child.Error("unknown command", slog.Int("cmdID", query.CmdID.Int()))

		return database.HandleCmdResult{}, command.ErrUnknownCmd
	}

	if err != nil {
		child.Error("failed to execute command", slog.String("err", err.Error()))

		return database.HandleCmdResult{}, err
	}

	if execRes.Value != "" {
		return database.HandleCmdResult{Result: execRes.Value.String()}, nil
	}

	return database.HandleCmdResult{Result: execRes.Result}, nil
}
