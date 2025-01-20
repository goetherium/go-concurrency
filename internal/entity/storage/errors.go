package storage

import (
	"errors"
)

var (
	ErrEmptyKey   = errors.New("empty key")
	ErrEmptyValue = errors.New("empty value")
)

// Nil reply returned by Storage when key does not exist.
const Nil = StorageError("nil") // nolint:errname

type StorageError string

func (e StorageError) Error() string {
	return string(e)
}
