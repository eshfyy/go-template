package infra

import (
	"context"
	"go-template/internal/domain"
)

type EventProducer interface {
	Publish(ctx context.Context, event domain.Event) error
}
