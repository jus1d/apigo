package log

import (
	"api/internal/config"
	"api/pkg/log/prettyslog"
	"log/slog"
	"os"
)

func InitDefault(env config.Env) {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = prettyslog.Init()
	case config.EnvDevelopment:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Key = "timestamp"
				}
				return a
			},
		}))
	case config.EnvProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Key = "timestamp"
				}
				return a
			},
		}))
	}

	logger = logger.With(
		slog.String("env", string(env)),
	)

	slog.SetDefault(logger)
}
