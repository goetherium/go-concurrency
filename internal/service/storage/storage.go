package storage

import (
	"log/slog"

	"database/internal/entity/command"
	"database/internal/entity/storage"
)

type engine interface {
	Set(*slog.Logger, command.Key, command.Value) (storage.ExecResult, error)
	Get(*slog.Logger, command.Key) (storage.ExecResult, error)
	Del(*slog.Logger, command.Key) (storage.ExecResult, error)
}

// todo: tests

type Storage struct {
	engine engine
}

func New(engine engine) *Storage {
	return &Storage{
		engine: engine,
	}
}

func (s Storage) Set(logger *slog.Logger, key command.Key, value command.Value) (storage.ExecResult, error) {
	return s.engine.Set(logger, key, value)
}

func (s Storage) Get(logger *slog.Logger, key command.Key) (storage.ExecResult, error) {
	return s.engine.Get(logger, key)
}

func (s Storage) Del(logger *slog.Logger, key command.Key) (storage.ExecResult, error) {
	return s.engine.Del(logger, key)
}
