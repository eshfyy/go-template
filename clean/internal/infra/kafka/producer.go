package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"go-template/internal/domain"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/twmb/franz-go/pkg/kgo"
)

var tracer = otel.Tracer("kafka")

type Producer struct {
	client *kgo.Client
	topic  string
	log    *zap.Logger
}

func NewProducer(brokers []string, topic string, log *zap.Logger) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(topic),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}
	return &Producer{client: client, topic: topic, log: log.Named("kafka_producer")}, nil
}

func (p *Producer) Publish(ctx context.Context, event domain.Event) error {
	ctx, span := tracer.Start(ctx, "kafka.publish",
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", p.topic),
			attribute.String("event.type", string(event.Type())),
		),
	)
	defer span.End()

	value, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("marshal event: %w", err)
	}

	record := &kgo.Record{
		Key:   []byte(event.Key()),
		Value: value,
		Headers: []kgo.RecordHeader{
			{Key: "event_type", Value: []byte(event.Type())},
		},
	}

	// Inject trace context into Kafka headers
	otel.GetTextMapPropagator().Inject(ctx, newHeaderCarrier(record))

	result := p.client.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		span.RecordError(err)
		p.log.Error("publish failed",
			zap.String("event_type", string(event.Type())),
			zap.Error(err),
		)
		return err
	}

	p.log.Debug("event published",
		zap.String("event_type", string(event.Type())),
		zap.String("key", event.Key()),
	)
	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}
