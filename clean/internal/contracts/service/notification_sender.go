package service

import (
	"context"

	"github.com/google/uuid"
)

type NotificationSenderService interface {
	Send(ctx context.Context, notificationID uuid.UUID) error
}
