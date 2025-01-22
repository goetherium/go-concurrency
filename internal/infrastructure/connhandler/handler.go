package connhandler

import (
	"errors"
	"fmt"
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
	cmdHandler     cmdHandler
	idleTimeout    time.Duration
	maxMessageSize int
}

func New(cmdHandler cmdHandler, idleTimeout time.Duration, maxMessageSize int) *Handler {
	return &Handler{
		cmdHandler:     cmdHandler,
		idleTimeout:    idleTimeout,
		maxMessageSize: maxMessageSize,
	}
}

func (h Handler) HandleConn(logger *slog.Logger, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("handle conn panic", slog.Any("err", err))
		}

		_ = conn.Close()
	}()

	timeout := time.Now().Add(h.idleTimeout)
	if err := conn.SetReadDeadline(timeout); err != nil {
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
	if len(buf) > h.maxMessageSize {
		return h.write(logger, conn, fmt.Sprintf("message exceeds max size %d bytes", h.maxMessageSize))
	}

	cmd := string(buf)
	cmd = strings.TrimRight(cmd, "\n")

	resp, err := h.cmdHandler.HandleCmd(logger, cmd)
	if err != nil {
		return h.write(logger, conn, err.Error())
	}

	return h.write(logger, conn, resp.String())
}

func (h Handler) write(logger *slog.Logger, conn net.Conn, data string) error {
	if _, err := conn.Write([]byte(data)); err != nil {
		logger.Error("failed to write to conn", slog.Any("err", err))

		return err
	}

	return nil
}
