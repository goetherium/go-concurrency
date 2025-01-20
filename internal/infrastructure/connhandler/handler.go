package connhandler

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

	"database/internal/entity/database"
	"database/internal/entity/iohelper"
)

type cmdHandler interface {
	HandleCmd(logger *slog.Logger, rowCmd string) (database.HandleCmdResult, error)
}

type Handler struct {
	cmdHandler cmdHandler
}

func New(cmdHandler cmdHandler) *Handler {
	return &Handler{
		cmdHandler: cmdHandler,
	}
}

func (h Handler) HandleConn(logger *slog.Logger, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("handle conn panic", slog.Any("err", err))
		}

		_ = conn.Close()
	}()

	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		logger.Error("set read deadline", slog.String("err", err.Error()))

		return
	}

	b := make([]byte, 0, 512)

	for {
		data, err := iohelper.Read(conn, b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info("client closed connection")
			} else {
				logger.Error("read conn", slog.Any("err", err))
			}

			return
		}

		logger.Info("read", slog.String("data", string(data)))

		if err = h.handleCmd(logger, conn, data); err != nil {
			return
		}
	}
}

func (h Handler) handleCmd(logger *slog.Logger, conn net.Conn, buf []byte) error {
	cmd := string(buf)
	cmd = strings.TrimRight(cmd, "\n")

	resp, err := h.cmdHandler.HandleCmd(logger, cmd)
	if err != nil {
		if _, err = conn.Write([]byte(err.Error())); err != nil {
			logger.Error("write to conn", slog.Any("err", err))

			return err
		}
	}

	if _, err = conn.Write([]byte(resp.String())); err != nil {
		logger.Error("write to conn", slog.Any("err", err))

		return err
	}

	return nil
}
