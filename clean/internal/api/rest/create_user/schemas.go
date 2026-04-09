package create_user

import "github.com/google/uuid"

type Request struct {
	Name       string  `json:"name" binding:"required"`
	Surname    *string `json:"surname"`
	TelegramID int64   `json:"telegram_id" binding:"required"`
}

type Response struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Surname    *string   `json:"surname,omitempty"`
	TelegramID int64     `json:"telegram_id"`
}
