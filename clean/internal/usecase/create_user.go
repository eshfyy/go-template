package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
)

type CreateUser struct {
	userRepo infra.UserRepository
}

func NewCreateUser(userRepo infra.UserRepository) *CreateUser {
	return &CreateUser{userRepo: userRepo}
}

func (u *CreateUser) Execute(ctx context.Context, input uc.CreateUserInput) (domain.User, error) {
	user, err := domain.NewUser(input.Name, input.Surname, input.TelegramID)
	if err != nil {
		return domain.User{}, err
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
