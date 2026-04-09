package usecase

import (
	"context"
	"go-template/internal/domain"
	"go-template/pkg/optional"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	Name       string
	Surname    optional.Optional[string]
	TelegramID int64
}

type CreateUser interface {
	Execute(ctx context.Context, input CreateUserInput) (domain.User, error)
}

type GetUser interface {
	Execute(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type ListUsersInput struct {
	Limit  int
	Offset int
}

type ListUsers interface {
	Execute(ctx context.Context, input ListUsersInput) ([]domain.User, error)
}

type UpdateUserInput struct {
	ID         uuid.UUID
	Name       string
	Surname    optional.Optional[string]
	TelegramID int64
}

type UpdateUser interface {
	Execute(ctx context.Context, input UpdateUserInput) (domain.User, error)
}

type DeleteUser interface {
	Execute(ctx context.Context, id uuid.UUID) error
}
