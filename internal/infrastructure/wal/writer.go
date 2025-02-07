package wal

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"database/internal/config"
	"database/internal/entity/logmodel"
)

const (
	firstFileIndex = "0"
	fileExt        = ".wal"
)

type writer struct {
	logger *slog.Logger

	dir          string
	fileMaxSize  uint64
	batchSize    int
	batchTimeout time.Duration

	filePath string
	file     *os.File
	lsn      Lsn

	tasks chan Task

	stopCtx context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func newWriter(logger *slog.Logger, cfg *config.Wal) *writer {
	return &writer{
		logger:       logger.With(logmodel.FieldModule, "wal_writer"),
		dir:          cfg.DataDir,
		fileMaxSize:  cfg.FileMaxSize,
		batchSize:    cfg.FlushBatchSize,
		batchTimeout: cfg.FlushBatchTimeout,

		tasks: make(chan Task),
	}
}

func (w *writer) start(filePath string, lsn Lsn) error {
	w.filePath = filePath
	w.lsn = lsn

	if err := w.openFile(); err != nil {
		return err
	}

	w.stopCtx, w.cancel = context.WithCancel(context.Background())

	w.wg.Add(1)

	go func() {
		defer w.wg.Done()

		w.worker()
	}()

	w.logger.Info("wal writer started")

	return nil
}

func (w *writer) stop() error {
	w.cancel()
	w.wg.Wait()

	close(w.tasks)

	if err := w.file.Close(); err != nil {
		w.logger.Error("failed to close wal file", "path", w.filePath, "err", err)

		return err
	}

	w.logger.Info("wal file closed")

	return nil
}

func (w *writer) submitTask(task Task) {
	w.tasks <- task
}

func (w *writer) worker() {
	ticker := time.NewTicker(w.batchTimeout)
	defer ticker.Stop()

	batch := make([]Task, 0, w.batchSize)

	for {
		select {
		case <-w.stopCtx.Done():
			w.logger.Info("stopping wal writer...")

			if len(batch) > 0 {
				w.handleBatch(batch)
			}

			return
		case task := <-w.tasks:
			batch = append(batch, task)

			if len(batch) >= w.batchSize {
				w.logger.Debug("handle wal batch by size")

				w.handleBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.logger.Debug("handle wal batch by ticker")

				w.handleBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

func (w *writer) handleBatch(batch []Task) {
	err := w.flush(batch)

	for _, task := range batch {
		task.Result <- err
	}
}

// Flush синхронно сбрасывает данные на диск.
func (w *writer) flush(batch []Task) error {
	var buf bytes.Buffer

	for _, task := range batch {
		w.lsn = w.lsn.inc()
		marshalCmd(&buf, w.lsn, task.Cmd)
	}

	if err := w.write(&buf); err != nil {
		return err
	}

	needNew, err := w.needNewFile()
	if err != nil {
		return err
	}

	if needNew {
		return w.createNewFile()
	}

	return nil
}

func (w *writer) write(buf *bytes.Buffer) error {
	_, err := w.file.Write(buf.Bytes())
	if err != nil {
		w.logger.Error("failed to write to wal file", "err", err)

		return err
	}

	if err = w.file.Sync(); err != nil {
		w.logger.Error("failed to sync wal wal file", "err", err)

		return err
	}

	return nil
}
