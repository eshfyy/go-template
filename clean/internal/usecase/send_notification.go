package usecase

import (
	"context"

	iservice "go-template/internal/contracts/service"

	"github.com/google/uuid"
)

type SendNotification struct {
	senderService iservice.NotificationSenderService
}

func NewSendNotification(senderService iservice.NotificationSenderService) *SendNotification {
	return &SendNotification{senderService: senderService}
}

func (u *SendNotification) Execute(ctx context.Context, notificationID uuid.UUID) error {
	return u.senderService.Send(ctx, notificationID)
}
