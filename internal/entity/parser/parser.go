package parser

import (
	"log/slog"
	"regexp"
	"strings"

	"database/internal/entity/command"
	"database/internal/entity/logmodel"
)

type Parser struct {
	r *regexp.Regexp
}

func New() *Parser {
	return &Parser{
		r: regexp.MustCompile("/(\\\\w+)/g"),
	}
}

func (p Parser) ParseCmd(logger *slog.Logger, rowCmd command.RawCmd) (command.Cmd, error) {
	child := logger.With(slog.String(logmodel.FieldAction, "ParseCmd"))

	tokens := strings.Fields(string(rowCmd))
	if len(tokens) == 0 {
		child.Error("invalid command", "cmd", rowCmd)

		return command.Cmd{}, command.ErrInvalidCmd
	}

	cmd := command.NewRawCmd(strings.ToUpper(tokens[0]))

	cmdID := cmd.ToCmdID()
	if cmdID == command.UnknownCommandID {
		child.Error("unknown command", "cmd", rowCmd)

		return command.Cmd{}, command.ErrUnknownCmd
	}

	gotCount := len(tokens) - 1
	wantCount := cmdID.ArgsCount()

	if gotCount != wantCount {
		child.Error("wrong args count",
			"cmd", rowCmd,
			slog.Int("got", gotCount),
			slog.Int("want", wantCount),
		)

		return command.Cmd{}, command.ErrInvalidArgsCount
	}

	args := command.NewArgs(tokens[1:])

	// Валидация аргументов
	for _, arg := range args {
		if p.r.MatchString(arg.String()) {
			child.Error("arg contains invalid symbols", "arg", arg)

			return command.Cmd{}, command.ErrInvalidArg
		}
	}

	child.Debug("cmd parse result", "cmd", cmd.String(), "args", args)

	return command.NewCmd(cmdID, args), nil
}
