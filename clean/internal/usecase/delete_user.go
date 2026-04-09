package usecase

import (
	"context"
	"go-template/internal/contracts/infra"

	"github.com/google/uuid"
)

type DeleteUser struct {
	userRepo infra.UserRepository
}

func NewDeleteUser(userRepo infra.UserRepository) *DeleteUser {
	return &DeleteUser{userRepo: userRepo}
}

func (u *DeleteUser) Execute(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}
