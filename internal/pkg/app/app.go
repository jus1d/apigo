package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/config"
	"api/pkg/log"

	httpserver "api/internal/transport/http"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{config}
}

func (a *App) Run() error {
	log.InitDefault(a.config.Env)

	server := httpserver.NewServer(a.config)

	errch := make(chan error, 1)
	go func() {
		errch <- server.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errch:
		return fmt.Errorf("server failed: %w", err)
	case <-quit:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
