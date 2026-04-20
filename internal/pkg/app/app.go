package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/config"
	"api/pkg/log"
	"api/pkg/log/sl"

	httpserver "api/internal/transport/http"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{config}
}

func (a *App) Run() {
	log.InitDefault(a.config.Env)

	server := httpserver.NewServer(a.config)

	go func() {
		if err := server.Run(); err != nil {
			slog.Error("failed to start server", sl.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", sl.Err(err))
		os.Exit(1)
	}
}
