package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	"go-template/internal/domain"

	"github.com/google/uuid"
)

type GetNotification struct {
	notificationRepo infra.NotificationRepository
}

func NewGetNotification(notificationRepo infra.NotificationRepository) *GetNotification {
	return &GetNotification{notificationRepo: notificationRepo}
}

func (u *GetNotification) Execute(ctx context.Context, id uuid.UUID) (domain.Notification, error) {
	return u.notificationRepo.GetByID(ctx, id)
}
