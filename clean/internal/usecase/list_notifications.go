package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
)

type ListNotifications struct {
	notificationRepo infra.NotificationRepository
}

func NewListNotifications(notificationRepo infra.NotificationRepository) *ListNotifications {
	return &ListNotifications{notificationRepo: notificationRepo}
}

func (u *ListNotifications) Execute(ctx context.Context, input uc.ListNotificationsInput) ([]domain.Notification, error) {
	return u.notificationRepo.ListByUserID(ctx, input.UserID, input.Limit, input.Offset)
}
