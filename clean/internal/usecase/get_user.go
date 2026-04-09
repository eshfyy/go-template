package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	"go-template/internal/domain"

	"github.com/google/uuid"
)

type GetUser struct {
	userRepo infra.UserRepository
}

func NewGetUser(userRepo infra.UserRepository) *GetUser {
	return &GetUser{userRepo: userRepo}
}

func (u *GetUser) Execute(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
