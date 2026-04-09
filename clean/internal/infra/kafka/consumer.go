package kafka

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler func(ctx context.Context, record *kgo.Record) error

type Consumer struct {
	client   *kgo.Client
	handlers map[string]Handler
	log      *zap.Logger
}

func NewConsumer(brokers []string, topic, groupID string, log *zap.Logger) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: %w", err)
	}
	return &Consumer{
		client:   client,
		handlers: make(map[string]Handler),
		log:      log.Named("kafka_consumer"),
	}, nil
}

func (c *Consumer) Register(eventType string, handler Handler) {
	c.handlers[eventType] = handler
	c.log.Info("registered handler", zap.String("event_type", eventType))
}

func (c *Consumer) Run(ctx context.Context) error {
	c.log.Info("polling started")
	for {
		fetches := c.client.PollFetches(ctx)
		if err := ctx.Err(); err != nil {
			return err
		}

		fetches.EachRecord(func(record *kgo.Record) {
			eventType := extractEventType(record)

			handler, ok := c.handlers[eventType]
			if !ok {
				c.log.Warn("no handler for event type", zap.String("type", eventType))
				return
			}

			// Extract trace context from Kafka headers to continue the trace
			propagator := otel.GetTextMapPropagator()
			recordCtx := propagator.Extract(ctx, newHeaderCarrier(record))

			recordCtx, span := tracer.Start(recordCtx, "kafka.consume",
				trace.WithAttributes(
					attribute.String("messaging.system", "kafka"),
					attribute.String("event.type", eventType),
					attribute.Int64("messaging.kafka.offset", record.Offset),
				),
				trace.WithSpanKind(trace.SpanKindConsumer),
			)

			if err := handler(recordCtx, record); err != nil {
				span.RecordError(err)
				c.log.Error("handle record failed",
					zap.String("type", eventType),
					zap.Int64("offset", record.Offset),
					zap.Error(err),
				)
			}
			span.End()
		})
	}
}

func extractEventType(record *kgo.Record) string {
	for _, h := range record.Headers {
		if h.Key == "event_type" {
			return string(h.Value)
		}
	}
	return ""
}

func (c *Consumer) Close() {
	c.client.Close()
}
