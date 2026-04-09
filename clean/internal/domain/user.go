package domain

import (
	"go-template/pkg/optional"
)

type User struct {
	BaseEntity
	Name       string
	Surname    optional.Optional[string]
	TelegramID int64
}

func NewUser(name string, surname optional.Optional[string], telegramID int64) (User, error) {
	fields := make(map[string]string)
	if name == "" {
		fields["name"] = "required"
	}
	if telegramID <= 0 {
		fields["telegram_id"] = "must be positive"
	}
	if len(fields) > 0 {
		return User{}, &ValidationError{Fields: fields}
	}

	return User{
		BaseEntity: NewBaseEntity(),
		Name:       name,
		Surname:    surname,
		TelegramID: telegramID,
	}, nil
}

func (u *User) UpdateProfile(name string, surname optional.Optional[string], telegramID int64) error {
	fields := make(map[string]string)
	if name == "" {
		fields["name"] = "required"
	}
	if telegramID <= 0 {
		fields["telegram_id"] = "must be positive"
	}
	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}

	u.Name = name
	u.Surname = surname
	u.TelegramID = telegramID
	return nil
}
