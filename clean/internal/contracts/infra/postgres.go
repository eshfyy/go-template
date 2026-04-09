package infra

import (
	"context"
	"go-template/internal/domain"
	"time"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Notification, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Notification, error)
	UpdateStatus(ctx context.Context, notification domain.Notification) error
	ListFailed(ctx context.Context, since time.Duration, limit int) ([]domain.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
