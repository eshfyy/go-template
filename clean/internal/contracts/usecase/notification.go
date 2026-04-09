package usecase

import (
	"context"
	"go-template/internal/domain"

	"github.com/google/uuid"
)

type CreateNotificationInput struct {
	UserID  uuid.UUID
	Title   string
	Text    string
	Channel domain.NotificationChannel
}

type CreateNotification interface {
	Execute(ctx context.Context, input CreateNotificationInput) (domain.Notification, error)
}

type GetNotification interface {
	Execute(ctx context.Context, id uuid.UUID) (domain.Notification, error)
}

type ListNotificationsInput struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListNotifications interface {
	Execute(ctx context.Context, input ListNotificationsInput) ([]domain.Notification, error)
}

type DeleteNotification interface {
	Execute(ctx context.Context, id uuid.UUID) error
}

type SendNotification interface {
	Execute(ctx context.Context, notificationID uuid.UUID) error
}

type RetryFailed interface {
	Execute(ctx context.Context) error
}
