package app

import (
	"go-template/internal/config"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(func() (*config.Config, error) {
		return config.Load("config")
	}),
	fx.Provide(func(cfg *config.Config) config.Postgres { return cfg.Postgres }),
	fx.Provide(func(cfg *config.Config) config.Kafka { return cfg.Kafka }),
	fx.Provide(func(cfg *config.Config) config.Redis { return cfg.Redis }),
	fx.Provide(func(cfg *config.Config) config.Telegram { return cfg.Telegram }),
	fx.Provide(func(cfg *config.Config) config.HTTP { return cfg.HTTP }),
	fx.Provide(func(cfg *config.Config) config.GRPC { return cfg.GRPC }),
	fx.Provide(func(cfg *config.Config) config.OTLP { return cfg.OTLP }),
)
