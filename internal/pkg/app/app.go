package app

import (
	"log/slog"

	"api/internal/config"
	"api/internal/version"
	"api/pkg/log"

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

	slog.Info("api: starting...", version.CommitAttr, version.BranchAttr)

	server := httpserver.NewServer(a.config)
	server.Run()
}
