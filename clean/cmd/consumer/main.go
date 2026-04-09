package main

import (
	"context"

	sendnotification "go-template/internal/api/consumer/send_notification"
	"go-template/internal/app"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
	"go-template/internal/infra/kafka"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		app.ConfigModule,
		app.LoggerModule,
		app.OTelModule("consumer"),
		app.PostgresModule,
		app.TelegramModule,
		app.ServiceModule,
		app.KafkaConsumerModule,
		app.UseCaseModule,

		fx.Invoke(startConsumer),
	).Run()
}

func startConsumer(lc fx.Lifecycle, consumer *kafka.Consumer, sendUC uc.SendNotification, log *zap.Logger) {
	consumer.Register(string(domain.NotificationCreated), sendnotification.New(sendUC))

	var cancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("starting consumer")
			var runCtx context.Context
			runCtx, cancel = context.WithCancel(context.Background())
			go func() {
				if err := consumer.Run(runCtx); err != nil && runCtx.Err() == nil {
					log.Error("consumer failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("stopping consumer")
			cancel()
			return nil
		},
	})
}
