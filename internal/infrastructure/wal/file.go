package wal

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (w *writer) openFile() error {
	var (
		file *os.File
		err  error
	)

	if w.filePath == "" {
		path := filepath.Join(w.dir, firstFileIndex+fileExt)
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else {
		file, err = os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}

	if err != nil {
		w.logger.Error("failed to open wal file", "path", w.filePath, "err", err)

		return err
	}

	w.file = file

	w.logger.Info("wal file opened", "path", filepath.Base(w.filePath))

	return nil
}

func (w *writer) needNewFile() (bool, error) {
	stat, err := w.file.Stat()
	if err != nil {
		w.logger.Error("failed to get wal file info", "err", err)

		return false, err
	}

	sizeExceeded := stat.Size() >= int64(w.fileMaxSize)

	return sizeExceeded, nil
}

func (w *writer) createNewFile() error {
	w.logger.Info("wal file size exceeds limit, switching to new one...")

	if err := w.file.Close(); err != nil {
		w.logger.Error("failed to close wal file", "err", err)

		return err
	}

	newName, err := w.newFileName()
	if err != nil {
		return err
	}

	w.filePath = filepath.Join(w.dir, newName)

	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		w.logger.Error("failed to create wal file", "path", w.filePath, "err", err)

		return err
	}

	w.file = file

	return nil
}

func (w *writer) newFileName() (string, error) {
	curName := filepath.Base(w.filePath)
	ext := filepath.Ext(curName)

	// Имя файла - число
	index, err := strconv.ParseUint(strings.TrimSuffix(curName, ext), 10, 64)
	if err != nil {
		w.logger.Error("failed to extract wal file index from name", "curName", curName, "err", err)

		return "", err
	}

	index++
	newName := strconv.FormatUint(index, 10) + fileExt

	w.logger.Info("new wal file name", "name", newName)

	return newName, nil
}
