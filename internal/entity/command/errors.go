package command

import (
	"errors"
)

var (
	ErrInvalidCmd       = errors.New("invalid command")
	ErrUnknownCmd       = errors.New("unknown command")
	ErrInvalidArgsCount = errors.New("invalid arguments count")
	ErrInvalidArg       = errors.New("invalid argument")
)
