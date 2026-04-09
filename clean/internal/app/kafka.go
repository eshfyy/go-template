package app

import (
	"context"
	"go-template/internal/config"
	"go-template/internal/contracts/infra"
	"go-template/internal/infra/kafka"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var KafkaProducerModule = fx.Module("kafka_producer",
	fx.Provide(
		fx.Annotate(
			func(lc fx.Lifecycle, cfg config.Kafka, log *zap.Logger) (*kafka.Producer, error) {
				p, err := kafka.NewProducer(cfg.Brokers, cfg.Topic, log)
				if err != nil {
					return nil, err
				}
				lc.Append(fx.Hook{
					OnStop: func(_ context.Context) error {
						p.Close()
						log.Info("kafka producer closed")
						return nil
					},
				})
				return p, nil
			},
			fx.As(new(infra.EventProducer)),
		),
	),
)

var KafkaConsumerModule = fx.Module("kafka_consumer",
	fx.Provide(func(lc fx.Lifecycle, cfg config.Kafka, log *zap.Logger) (*kafka.Consumer, error) {
		c, err := kafka.NewConsumer(cfg.Brokers, cfg.Topic, cfg.GroupID, log)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				c.Close()
				log.Info("kafka consumer closed")
				return nil
			},
		})
		return c, nil
	}),
)
