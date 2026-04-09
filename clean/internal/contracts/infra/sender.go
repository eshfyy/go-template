package infra

import (
	"context"
	"go-template/internal/domain"
)

type SendRequest struct {
	Notification domain.Notification
	Recipient    domain.User
}

type NotificationSender interface {
	Send(ctx context.Context, req SendRequest) error
	Channel() domain.NotificationChannel
}
