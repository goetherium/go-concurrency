package storage

import (
	"context"
	"log/slog"

	"database/internal/entity/command"
	"database/internal/entity/logmodel"
	"database/internal/entity/storage"
)

type walLayer interface {
	Open(ctx context.Context, applyCmd func(cmd command.Cmd) error) error
	Shutdown() error
	WriteCmd(logger *slog.Logger, cmd command.Cmd) error
}

type engineLayer interface {
	Set(*slog.Logger, command.Key, command.Value) (storage.ExecResult, error)
	Get(*slog.Logger, command.Key) (storage.ExecResult, error)
	Del(*slog.Logger, command.Key) (storage.ExecResult, error)
}

type Storage struct {
	logger *slog.Logger
	wal    walLayer
	engine engineLayer
}

func New(logger *slog.Logger, wal walLayer, engine engineLayer) *Storage {
	return &Storage{
		logger: logger,
		wal:    wal,
		engine: engine,
	}
}

func (s Storage) Open(ctx context.Context) error {
	return s.wal.Open(ctx, s.applyWalCmd(s.logger))
}

func (s Storage) Shutdown() error {
	return s.wal.Shutdown()
}

func (s Storage) ExecCmd(logger *slog.Logger, cmd command.Cmd) (storage.ExecResult, error) {
	child := logger.With(slog.String(logmodel.FieldAction, "ExecCmd"))

	if cmd.ID().IsWalCmd() {
		if err := s.wal.WriteCmd(child, cmd); err != nil {
			return storage.ExecResult{}, err
		}
	}

	return s.applyToEngine(child, cmd)
}

func (s Storage) applyToEngine(logger *slog.Logger, cmd command.Cmd) (storage.ExecResult, error) {
	logger.Debug("applying cmd to engine", "cmdID", cmd.ID(), "args", cmd.Args())

	var (
		args = cmd.Args()
		key  = command.NewKey(args[0])
		res  storage.ExecResult
		err  error
	)

	switch cmd.ID() {
	case command.SetCommandID:
		res, err = s.engine.Set(logger, key, command.NewValue(args[1]))
	case command.GetCommandID:
		res, err = s.engine.Get(logger, key)
	case command.DelCommandID:
		res, err = s.engine.Del(logger, key)
	default:
		logger.Error("unknown command", "cmdID", cmd.ID())

		return storage.ExecResult{}, command.ErrUnknownCmd
	}

	if err != nil {
		logger.Error("failed to execute command", slog.String("err", err.Error()))

		return storage.ExecResult{}, err
	}

	return res, nil
}

func (s Storage) applyWalCmd(logger *slog.Logger) func(cmd command.Cmd) error {
	return func(cmd command.Cmd) error {
		if !cmd.ID().IsWalCmd() {
			return nil
		}

		_, err := s.applyToEngine(logger, cmd)

		return err
	}
}
