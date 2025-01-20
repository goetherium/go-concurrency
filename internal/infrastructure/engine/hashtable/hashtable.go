package hashtable

import (
	"log/slog"
	"sync"

	"database/internal/entity/command"
	"database/internal/entity/storage"
)

// todo: tests

const resultOK = "OK"

type Engine struct {
	data map[command.Key]command.Value
	mx   sync.RWMutex
}

func New(initSize int) *Engine {
	return &Engine{
		data: make(map[command.Key]command.Value, initSize),
	}
}

func (e *Engine) Set(logger *slog.Logger, key command.Key, value command.Value) (storage.ExecResult, error) {
	if key == "" {
		logger.Error("set key is empty")

		return storage.ExecResult{}, storage.ErrEmptyKey
	}

	if value == "" {
		logger.Error("set value is empty")

		return storage.ExecResult{}, storage.ErrEmptyValue
	}

	e.mx.Lock()
	defer e.mx.Unlock()

	e.data[key] = value

	res := resultOK
	logger.Debug("SET", slog.String("result", res))

	return storage.ExecResult{Result: res}, nil
}

func (e *Engine) Get(logger *slog.Logger, key command.Key) (storage.ExecResult, error) {
	if key == "" {
		logger.Error("get key is empty")

		return storage.ExecResult{}, storage.ErrEmptyKey
	}

	e.mx.RLock()
	defer e.mx.RUnlock()

	value, ok := e.data[key]
	if !ok {
		logger.Error("key does not exist", slog.Any("key", key))

		return storage.ExecResult{}, storage.Nil
	}

	logger.Debug("GET", slog.String("result", value.String()))

	return storage.ExecResult{Value: value}, nil
}

func (e *Engine) Del(logger *slog.Logger, key command.Key) (storage.ExecResult, error) {
	if key == "" {
		logger.Error("delete key is empty")

		return storage.ExecResult{}, storage.ErrEmptyKey
	}

	e.mx.Lock()
	defer e.mx.Unlock()

	delete(e.data, key)

	res := resultOK
	logger.Debug("DEL", slog.String("result", res))

	return storage.ExecResult{Result: res}, nil
}
