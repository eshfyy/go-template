package entities

import (
	vo "go-template/internal/domain/valueobjects"
	op "go-template/pkg/optional"
)

type User struct {
	BaseEntity
	Name    string
	Surname op.Optional[string]
	Email   vo.Email
}

func NewUser(name string, surname op.Optional[string], email string) (User, error) {
	emailVo, err := vo.NewEmail(email)
	if err != nil {
		return User{}, err
	}
	return User{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Surname:    surname,
		Email:      emailVo,
	}, nil
}
