package wal

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"database/internal/config"
	"database/internal/entity/command"
	"database/internal/entity/logmodel"
)

type cmdParser interface {
	ParseCmd(logger *slog.Logger, rowCmd command.RawCmd) (command.Cmd, error)
}

type Wal struct {
	logger *slog.Logger
	writer *writer
	parser cmdParser
	dir    string
}

func NewWal(logger *slog.Logger, cfg *config.Wal, parser cmdParser) *Wal {
	return &Wal{
		logger: logger,
		writer: newWriter(logger, cfg),
		parser: parser,
		dir:    cfg.DataDir,
	}
}

func (w *Wal) Open(ctx context.Context, applyCmd func(cmd command.Cmd) error) error {
	filePath, lsn, err := w.loadAll(ctx, applyCmd)
	if err != nil {
		return err
	}

	return w.writer.start(filePath, lsn)
}

func (w *Wal) Shutdown() error {
	return w.writer.stop()
}

func (w *Wal) WriteCmd(logger *slog.Logger, cmd command.Cmd) error {
	logger.Debug("writing cmd in wal", "cmdID", cmd.ID(), "args", cmd.Args())

	if err := w.writeCmd(cmd); err != nil {
		logger.Error("failed to write cmd into wal", "error", err)

		return err
	}

	return nil
}

func (w *Wal) loadAll(ctx context.Context, applyCmd func(cmd command.Cmd) error) (string, Lsn, error) {
	child := w.logger.With(slog.String(logmodel.FieldAction, "load_wals"))

	files, err := os.ReadDir(w.dir)
	if err != nil {
		child.Error("failed to read wal dir", "err", err)

		return "", 0, err
	}

	var (
		path string
		lsn  Lsn
	)

	for _, file := range files {
		if ctx.Err() != nil {
			return "", 0, ctx.Err()
		}

		path = filepath.Join(w.dir, file.Name())
		if lsn, err = w.loadOne(ctx, child, path, applyCmd); err != nil {
			return "", 0, err
		}
	}

	// Поскольку os.ReadDir возвращает файлы в каталоге, отсортированные по имени,
	// последний файл и будет текущим.
	return path, lsn, nil
}

func (w *Wal) loadOne(
	ctx context.Context,
	logger *slog.Logger,
	path string,
	applyCmd func(cmd command.Cmd) error,
) (Lsn, error) {
	child := logger.With(slog.String(logmodel.FieldAction, "load_file"), "name", filepath.Base(path))
	child.Info("loading wal file...")

	file, err := os.Open(path)
	if err != nil {
		child.Error("failed to open wal file", "path", path, "err", err)

		return 0, err
	}

	defer func() {
		_ = file.Close()
	}()

	// Применим все команды из файла
	scanner := bufio.NewScanner(file)
	var lastLsn Lsn

	for scanner.Scan() {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		line := scanner.Text()

		lsn, rawCmd, err := parseLine(line)
		if err != nil {
			child.Error("failed to parse line", "line", line, "err", err)

			continue
		}

		lastLsn = lsn

		child.Debug("parsed line", "lsn", lsn, "rawCmd", rawCmd)

		cmd, err := w.parser.ParseCmd(child, rawCmd)
		if err != nil {
			child.Error("failed to parse cmd", "rawCmd", rawCmd, "err", err)

			continue
		}

		child.Debug("applying action to wal line", "lsn", lsn)

		if err := applyCmd(cmd); err != nil {
			child.Error("failed to apply wal command", "cmd", cmd, "err", err)
		}
	}

	if err := scanner.Err(); err != nil {
		child.Error("failed to scan wal file", "err", err)

		return 0, err
	}

	child.Info("wal file loaded successfully", "lastLSN", lastLsn)

	return lastLsn, nil
}

func (w *Wal) writeCmd(cmd command.Cmd) error {
	result := make(chan error)

	task := Task{
		Cmd:    cmd,
		Result: result,
	}

	w.writer.submitTask(task)

	return <-result
}
