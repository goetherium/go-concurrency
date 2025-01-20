package server

import (
	"log/slog"
	"net"
)

type connHandler interface {
	HandleConn(logger *slog.Logger, conn net.Conn)
}

type Server struct {
	logger  *slog.Logger
	handler connHandler
}

func New(logger *slog.Logger, handler connHandler) *Server {
	return &Server{
		logger:  logger,
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

		// todo: acquire semaphore
		go s.handler.HandleConn(s.logger, conn)
	}

}
