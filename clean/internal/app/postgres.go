package app

import (
	"context"
	"go-template/internal/config"
	"go-template/internal/contracts/infra"
	"go-template/internal/infra/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var PostgresModule = fx.Module("postgres",
	fx.Provide(func(lc fx.Lifecycle, cfg config.Postgres, log *zap.Logger) (*pgxpool.Pool, error) {
		pool, err := postgres.New(cfg, log)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				// ctx has fx startup timeout — ping won't hang forever
				if err := pool.Ping(ctx); err != nil {
					return err
				}
				log.Info("postgres connected")
				return nil
			},
			OnStop: func(_ context.Context) error {
				pool.Close()
				log.Info("postgres connection closed")
				return nil
			},
		})
		return pool, nil
	}),
	fx.Provide(
		fx.Annotate(
			postgres.NewNotificationRepository,
			fx.As(new(infra.NotificationRepository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			postgres.NewUserRepository,
			fx.As(new(infra.UserRepository)),
		),
	),
)
