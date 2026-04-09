package app

import (
	"context"
	"go-template/internal/config"
	pkgotel "go-template/pkg/otel"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func OTelModule(serviceSuffix string) fx.Option {
	return fx.Module("otel",
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					shutdown, err := pkgotel.Init(ctx, pkgotel.Config{
						ServiceName: cfg.App.Name + "-" + serviceSuffix,
						Endpoint:    cfg.OTLP.Endpoint,
					})
					if err != nil {
						return err
					}
					lc.Append(fx.Hook{
						OnStop: func(ctx context.Context) error {
							log.Info("shutting down otel")
							return shutdown(ctx)
						},
					})
					return nil
				},
			})
		}),
	)
}
