package main

import (
	"context"
	"fmt"

	retryfailed "go-template/internal/api/worker/retry_failed"
	"go-template/internal/app"
	"go-template/internal/config"
	uc "go-template/internal/contracts/usecase"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		app.ConfigModule,
		app.LoggerModule,
		app.OTelModule("worker"),
		app.PostgresModule,
		app.TelegramModule,
		app.ServiceModule,
		app.UseCaseModule,

		fx.Provide(newAsynqServer),
		fx.Provide(newAsynqScheduler),

		fx.Invoke(startWorker),
	).Run()
}

func newAsynqServer(cfg config.Redis) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{Concurrency: 1},
	)
}

func newAsynqScheduler(cfg config.Redis) *asynq.Scheduler {
	return asynq.NewScheduler(asynq.RedisClientOpt{Addr: cfg.Addr}, nil)
}

func startWorker(
	lc fx.Lifecycle,
	srv *asynq.Server,
	scheduler *asynq.Scheduler,
	retryUC uc.RetryFailed,
	log *zap.Logger,
) error {
	if _, err := scheduler.Register("@every 1m", retryfailed.NewTask()); err != nil {
		return fmt.Errorf("register scheduler task: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("starting worker")

			go func() {
				if err := scheduler.Run(); err != nil {
					log.Error("scheduler failed", zap.Error(err))
				}
			}()

			mux := asynq.NewServeMux()
			mux.HandleFunc(retryfailed.TaskType, retryfailed.New(retryUC))

			go func() {
				if err := srv.Run(mux); err != nil {
					log.Error("asynq server failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			scheduler.Shutdown()
			srv.Shutdown()
			log.Info("worker stopped")
			return nil
		},
	})

	return nil
}
