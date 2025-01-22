package server

import (
	"log/slog"
	"net"

	"database/internal/infrastructure/semaphore"
)

type connHandler interface {
	HandleConn(logger *slog.Logger, conn net.Conn)
}

type Server struct {
	logger  *slog.Logger
	sema    *semaphore.Semaphore
	handler connHandler
}

func New(logger *slog.Logger, sema *semaphore.Semaphore, handler connHandler) *Server {
	return &Server{
		logger:  logger,
		sema:    sema,
		handler: handler,
	}
}

func (s Server) ListenAndServe(addr string) error {
	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		return listenErr
	}

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("accept error", slog.String("err", err.Error()))

			continue
		}

		s.sema.Acquire()

		go func() {
			defer s.sema.Release()

			s.handler.HandleConn(s.logger, conn)
		}()
	}

}
