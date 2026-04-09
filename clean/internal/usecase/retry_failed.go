package usecase

import (
	"context"
	"fmt"
	"go-template/internal/contracts/infra"
	iservice "go-template/internal/contracts/service"
	"time"

	"go.uber.org/zap"
)

type RetryFailed struct {
	notificationRepo infra.NotificationRepository
	senderService    iservice.NotificationSenderService
	log              *zap.Logger
}

func NewRetryFailed(
	notificationRepo infra.NotificationRepository,
	senderService iservice.NotificationSenderService,
	log *zap.Logger,
) *RetryFailed {
	return &RetryFailed{
		notificationRepo: notificationRepo,
		senderService:    senderService,
		log:              log.Named("retry_failed"),
	}
}

func (uc *RetryFailed) Execute(ctx context.Context) error {
	notifications, err := uc.notificationRepo.ListFailed(ctx, 30*time.Minute, 100)
	if err != nil {
		return fmt.Errorf("list failed: %w", err)
	}

	if len(notifications) == 0 {
		uc.log.Debug("no failed notifications to retry")
		return nil
	}

	uc.log.Info("retrying failed notifications", zap.Int("count", len(notifications)))

	for _, n := range notifications {
		if err := uc.senderService.Send(ctx, n.ID); err != nil {
			uc.log.Error("retry failed",
				zap.String("notification_id", n.ID.String()),
				zap.Error(err),
			)
			continue
		}
		uc.log.Info("retry succeeded", zap.String("notification_id", n.ID.String()))
	}

	return nil
}
