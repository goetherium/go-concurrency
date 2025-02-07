package wal

import (
	"database/internal/entity/command"
)

type Lsn uint64 // Last sequence number - последний номер записи в WAL.

func newLsn(in uint64) Lsn {
	return Lsn(in)
}

func (lsn Lsn) inc() Lsn {
	return lsn + 1
}

type Task struct {
	Cmd    command.Cmd
	Result chan error
}
