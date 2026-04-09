package app

import (
	"go-template/internal/config"
	"go-template/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var LoggerModule = fx.Module("logger",
	fx.Provide(func(cfg *config.Config) (*zap.Logger, error) {
		return logger.New(cfg.App.LogLevel, string(cfg.App.Env))
	}),
)
