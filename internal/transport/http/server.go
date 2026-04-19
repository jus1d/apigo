package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api/internal/config"
	v1 "api/internal/transport/http/v1"
	"api/pkg/log/sl"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(c *config.Config) *Server {
	router := v1.NewRouter(c)

	return &Server{
		httpServer: &http.Server{
			Addr:         c.Server.Address,
			Handler:      router.InitRoutes(),
			ReadTimeout:  c.Server.Timeout,
			WriteTimeout: c.Server.Timeout,
			IdleTimeout:  c.Server.IdleTimeout,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to start server", sl.Err(err))
				os.Exit(1)
			}
		}
	}()

	slog.Info("api: started", slog.String("address", s.httpServer.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	slog.Info("api: shutting down...")

	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		slog.Error("api: error occurred on server shutting down", sl.Err(err))
		os.Exit(1)
	}

	slog.Info("api: server stopped")
}
