package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	_defaultReadTimeout     = 30 * time.Second
	_defaultWriteTimeout    = 30 * time.Second
	_defaultShutdownTimeout = 30 * time.Second
	_defaultAddr            = ":80"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler *gin.Engine, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}
	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func(s *Server) {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}(s)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
