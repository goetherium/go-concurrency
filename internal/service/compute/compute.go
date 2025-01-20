package compute

import (
	"log/slog"

	"database/internal/entity/command"
)

type cmdParser interface {
	ParseCmd(logger *slog.Logger, rowCmd string) (command.Query, error)
}

type Compute struct {
	p cmdParser
}

func New(p cmdParser) *Compute {
	return &Compute{
		p: p,
	}
}

func (h Compute) ParseCmd(logger *slog.Logger, rowCmd string) (command.Query, error) {
	return h.p.ParseCmd(logger, rowCmd)
}
