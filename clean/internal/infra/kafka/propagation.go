package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/propagation"
)

// headerCarrier adapts kgo.Record headers to OTel TextMapCarrier.
type headerCarrier struct {
	record *kgo.Record
}

func newHeaderCarrier(record *kgo.Record) *headerCarrier {
	return &headerCarrier{record: record}
}

func (c *headerCarrier) Get(key string) string {
	for _, h := range c.record.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *headerCarrier) Set(key, value string) {
	for i, h := range c.record.Headers {
		if h.Key == key {
			c.record.Headers[i].Value = []byte(value)
			return
		}
	}
	c.record.Headers = append(c.record.Headers, kgo.RecordHeader{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *headerCarrier) Keys() []string {
	keys := make([]string, len(c.record.Headers))
	for i, h := range c.record.Headers {
		keys[i] = h.Key
	}
	return keys
}

var _ propagation.TextMapCarrier = (*headerCarrier)(nil)
