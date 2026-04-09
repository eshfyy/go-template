package usecase

import (
	"context"
	"fmt"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"

	"go.uber.org/zap"
)

type CreateNotification struct {
	notificationRepo infra.NotificationRepository
	eventProducer    infra.EventProducer
	log              *zap.Logger
}

func NewCreateNotification(
	notificationRepo infra.NotificationRepository,
	eventProducer infra.EventProducer,
	log *zap.Logger,
) *CreateNotification {
	return &CreateNotification{
		notificationRepo: notificationRepo,
		eventProducer:    eventProducer,
		log:              log.Named("create_notification"),
	}
}

func (u *CreateNotification) Execute(ctx context.Context, input uc.CreateNotificationInput) (domain.Notification, error) {
	notification, err := domain.NewNotification(input.UserID, input.Title, input.Text, input.Channel)
	if err != nil {
		return domain.Notification{}, err
	}

	if err := u.notificationRepo.Create(ctx, notification); err != nil {
		return domain.Notification{}, fmt.Errorf("create notification: %w", err)
	}

	event := domain.NewNotificationCreatedEvent(notification.ID)
	if err := u.eventProducer.Publish(ctx, event); err != nil {
		u.log.Error("failed to publish event",
			zap.String("notification_id", notification.ID.String()),
			zap.Error(err),
		)
		return domain.Notification{}, fmt.Errorf("publish event: %w", err)
	}

	u.log.Info("notification created",
		zap.String("id", notification.ID.String()),
		zap.String("channel", string(notification.Channel)),
	)

	return notification, nil
}
