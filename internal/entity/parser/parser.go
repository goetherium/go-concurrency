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

func (p Parser) ParseCmd(logger *slog.Logger, rowCmd string) (command.Query, error) {
	child := logger.With(slog.String(logmodel.FieldAction, "ParseCmd"))

	tokens := strings.Fields(rowCmd)
	if len(tokens) == 0 {
		child.Error("invalid command", slog.String("cmd", rowCmd))

		return command.Query{}, command.ErrInvalidCmd
	}

	cmd := command.NewRawCmd(strings.ToUpper(tokens[0]))

	cmdID := cmd.ToCmdID()
	if cmdID == command.UnknownCommandID {
		child.Error("unknown command", slog.String("cmd", rowCmd))

		return command.Query{}, command.ErrUnknownCmd
	}

	gotCount := len(tokens) - 1
	wantCount := cmdID.ArgsCount()

	if gotCount != wantCount {
		child.Error("wrong args count",
			slog.String("cmd", rowCmd),
			slog.Int("got", gotCount),
			slog.Int("want", wantCount),
		)

		return command.Query{}, command.ErrInvalidArgsCount
	}

	args := tokens[1:]

	// Валидация аргументов
	for _, arg := range args {
		if p.r.MatchString(arg) {
			child.Error("arg contains invalid symbols", slog.String("arg", arg))

			return command.Query{}, command.ErrInvalidArg
		}
	}

	child.Debug("cmd parse result",
		slog.String("cmd", cmd.String()),
		slog.String("args", strings.Join(args, " ")),
	)

	return command.Query{
		CmdID: cmdID,
		Args:  args,
	}, nil
}
