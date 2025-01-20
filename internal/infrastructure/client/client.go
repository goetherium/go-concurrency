package client

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"time"

	"database/internal/entity/iohelper"
)

func Run(logger *slog.Logger, hostPort string) {
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		logger.Error("failed to connect to host", "err", err)

		return
	}

	defer func() {
		_ = conn.Close()
	}()

	if err = conn.SetWriteDeadline(time.Time{}); err != nil {
		logger.Error("failed to set write deadline", "err", err)

		return
	}

	if err = conn.SetReadDeadline(time.Time{}); err != nil {
		logger.Error("failed to set read deadline", "err", err)

		return
	}

	writeBuf := make([]byte, 512)
	readBuf := make([]byte, 0, 512)

	for {
		n, err := os.Stdin.Read(writeBuf)
		if err != nil {
			logger.Error("failed to read from stdin", "err", err)

			return
		}

		if _, err = conn.Write(writeBuf[:n]); err != nil {
			logger.Error("failed to send data to server", "err", err)

			return
		}

		response, err := iohelper.Read(conn, readBuf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info("server closed connection")
			} else {
				logger.Error("read", slog.Any("err", err))
			}

			return
		}

		response = append(response, '\n')
		_, _ = os.Stdout.Write(response)
	}
}
