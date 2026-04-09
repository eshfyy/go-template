package postgres

import (
	"context"
	"fmt"
	"go-template/internal/domain"
	"go-template/internal/infra/postgres/sqlc"
	"go-template/pkg/optional"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	q *sqlc.Queries
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{q: sqlc.New(pool)}
}

func (r *NotificationRepository) Create(ctx context.Context, n domain.Notification) error {
	err := r.q.CreateNotification(ctx, sqlc.CreateNotificationParams{
		ID:        uuidToPg(n.ID),
		UserID:    uuidToPg(n.UserID),
		Title:     n.Title,
		Text:      n.Text,
		Channel:   string(n.Channel),
		Status:    string(n.Status),
		CreatedAt: timeToPgTimestamptz(n.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("create notification: %w", mapError(err, "notification"))
	}
	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Notification, error) {
	row, err := r.q.GetNotificationByID(ctx, uuidToPg(id))
	if err != nil {
		return domain.Notification{}, fmt.Errorf("get notification %s: %w", id, mapError(err, "notification"))
	}
	return notificationFromRow(row), nil
}

func (r *NotificationRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error) {
	rows, err := r.q.ListNotificationsByUserID(ctx, sqlc.ListNotificationsByUserIDParams{
		UserID: uuidToPg(userID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list notifications for user %s: %w", userID, err)
	}
	notifications := make([]domain.Notification, len(rows))
	for i, row := range rows {
		notifications[i] = notificationFromRow(row)
	}
	return notifications, nil
}

func (r *NotificationRepository) UpdateStatus(ctx context.Context, n domain.Notification) error {
	params := sqlc.UpdateNotificationStatusParams{
		ID:     uuidToPg(n.ID),
		Status: string(n.Status),
	}
	if sentAt, ok := n.SentAt.Get(); ok {
		params.SentAt = timeToPgTimestamptz(sentAt)
	}
	err := r.q.UpdateNotificationStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("update notification status %s: %w", n.ID, mapError(err, "notification"))
	}
	return nil
}

func (r *NotificationRepository) ListFailed(ctx context.Context, since time.Duration, limit int) ([]domain.Notification, error) {
	rows, err := r.q.ListFailedNotifications(ctx, sqlc.ListFailedNotificationsParams{
		Limit: int32(limit),
		Since: durationToPgInterval(since),
	})
	if err != nil {
		return nil, fmt.Errorf("list failed notifications: %w", err)
	}
	notifications := make([]domain.Notification, len(rows))
	for i, row := range rows {
		notifications[i] = notificationFromRow(row)
	}
	return notifications, nil
}

func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteNotification(ctx, uuidToPg(id))
	if err != nil {
		return fmt.Errorf("delete notification %s: %w", id, mapError(err, "notification"))
	}
	return nil
}

// notificationFromRow hydrates a domain.Notification from a DB row.
// Validation is not called — data in the database is considered valid.
func notificationFromRow(row sqlc.Notification) domain.Notification {
	n := domain.Notification{
		BaseEntity: domain.BaseEntity{
			ID:        pgToUUID(row.ID),
			CreatedAt: row.CreatedAt.Time,
		},
		UserID:  pgToUUID(row.UserID),
		Title:   row.Title,
		Text:    row.Text,
		Channel: domain.NotificationChannel(row.Channel),
		Status:  domain.NotificationStatus(row.Status),
	}
	if row.SentAt.Valid {
		n.SentAt = optional.Some(row.SentAt.Time)
	}
	return n
}
