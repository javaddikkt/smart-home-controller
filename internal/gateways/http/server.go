package http

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/gateways/http/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host       string
	port       uint16
	router     *gin.Engine
	wsHandler  *WebSocketHandler
	httpServer *http.Server
}

func NewServer(useCases types.UseCases, options ...func(*Server)) *Server {
	r := gin.Default()
	ws := NewWebSocketHandler(useCases)
	setupRouter(r, useCases, ws)

	s := &Server{router: r, host: "localhost", port: 8080, wsHandler: ws}
	for _, o := range options {
		o(s)
	}

	return s
}

func WithHost(host string) func(*Server) {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port uint16) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		shutdownErr := s.shutdownServer()
		return shutdownErr
	case err := <-errChan:
		return err
	}
}

func (s *Server) shutdownServer() error {
	if s.wsHandler != nil {
		_ = s.wsHandler.Shutdown()
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
