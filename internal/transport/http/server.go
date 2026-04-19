package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"api/internal/config"
	"api/internal/transport/http/middleware"
	v1 "api/internal/transport/http/v1"
	"api/pkg/requestid"
	"api/pkg/requestlog"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(c *config.Config) *Server {
	e := echo.New()

	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(requestid.New)
	e.Use(requestlog.Completed)
	e.Pre(echomw.RemoveTrailingSlash())

	switch c.Env {
	case config.EnvLocal, config.EnvDevelopment:
		e.Use(middleware.CORS)
	}

	api := e.Group("/api")

	v1Group := api.Group("/v1")
	v1Router := v1.NewRouter()
	v1Router.Register(v1Group)

	return &Server{
		httpServer: &http.Server{
			Addr:         c.Server.Address,
			Handler:      e,
			ReadTimeout:  c.Server.Timeout,
			WriteTimeout: c.Server.Timeout,
			IdleTimeout:  c.Server.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	slog.Info("api: started", slog.String("address", s.httpServer.Addr))

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("api: shutting down...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	slog.Info("api: server stopped")
	return nil
}
