package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
)

type ListUsers struct {
	userRepo infra.UserRepository
}

func NewListUsers(userRepo infra.UserRepository) *ListUsers {
	return &ListUsers{userRepo: userRepo}
}

func (u *ListUsers) Execute(ctx context.Context, input uc.ListUsersInput) ([]domain.User, error) {
	return u.userRepo.List(ctx, input.Limit, input.Offset)
}
