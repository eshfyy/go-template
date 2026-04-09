package service

import (
	"context"
	"fmt"
	"go-template/internal/contracts/infra"
	"go-template/internal/domain"
	"go-template/pkg/logger"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("service")

type NotificationSenderService struct {
	notificationRepo infra.NotificationRepository
	userRepo         infra.UserRepository
	senders          map[domain.NotificationChannel]infra.NotificationSender
	log              *zap.Logger
}

func NewNotificationSenderService(
	notificationRepo infra.NotificationRepository,
	userRepo infra.UserRepository,
	senders []infra.NotificationSender,
	log *zap.Logger,
) *NotificationSenderService {
	senderMap := make(map[domain.NotificationChannel]infra.NotificationSender, len(senders))
	for _, s := range senders {
		senderMap[s.Channel()] = s
	}
	return &NotificationSenderService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		senders:          senderMap,
		log:              log.Named("notification_sender"),
	}
}

func (s *NotificationSenderService) Send(ctx context.Context, notificationID uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "notification_sender.send",
		trace.WithAttributes(attribute.String("notification_id", notificationID.String())),
	)
	defer span.End()

	log := logger.WithTrace(ctx, s.log)

	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("get notification: %w", err)
	}

	span.SetAttributes(
		attribute.String("channel", string(notification.Channel)),
		attribute.String("user_id", notification.UserID.String()),
	)

	sender, ok := s.senders[notification.Channel]
	if !ok {
		return fmt.Errorf("unsupported channel: %s", notification.Channel)
	}

	user, err := s.userRepo.GetByID(ctx, notification.UserID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("get user: %w", err)
	}

	log.Info("sending notification",
		zap.String("notification_id", notificationID.String()),
		zap.String("channel", string(notification.Channel)),
		zap.String("user_id", notification.UserID.String()),
	)

	err = sender.Send(ctx, infra.SendRequest{
		Notification: notification,
		Recipient:    user,
	})
	if err != nil {
		notification.MarkFailed()
		if updateErr := s.notificationRepo.UpdateStatus(ctx, notification); updateErr != nil {
			span.RecordError(updateErr)
			log.Error("failed to update status after send failure",
				zap.String("notification_id", notificationID.String()),
				zap.Error(updateErr),
			)
			return fmt.Errorf("send failed: %w, status update also failed: %v", err, updateErr)
		}
		span.RecordError(err)
		log.Error("send failed",
			zap.String("notification_id", notificationID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("send notification: %w", err)
	}

	notification.MarkSent()
	if err := s.notificationRepo.UpdateStatus(ctx, notification); err != nil {
		span.RecordError(err)
		return fmt.Errorf("update status: %w", err)
	}

	log.Info("notification sent",
		zap.String("notification_id", notificationID.String()),
	)

	return nil
}
