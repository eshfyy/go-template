package usecase

import (
	"context"
	"go-template/internal/contracts/infra"

	"github.com/google/uuid"
)

type DeleteNotification struct {
	notificationRepo infra.NotificationRepository
}

func NewDeleteNotification(notificationRepo infra.NotificationRepository) *DeleteNotification {
	return &DeleteNotification{notificationRepo: notificationRepo}
}

func (u *DeleteNotification) Execute(ctx context.Context, id uuid.UUID) error {
	return u.notificationRepo.Delete(ctx, id)
}
