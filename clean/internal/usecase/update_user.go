package usecase

import (
	"context"
	"go-template/internal/contracts/infra"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"
)

type UpdateUser struct {
	userRepo infra.UserRepository
}

func NewUpdateUser(userRepo infra.UserRepository) *UpdateUser {
	return &UpdateUser{userRepo: userRepo}
}

func (u *UpdateUser) Execute(ctx context.Context, input uc.UpdateUserInput) (domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.User{}, err
	}

	if err := user.UpdateProfile(input.Name, input.Surname, input.TelegramID); err != nil {
		return domain.User{}, err
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
