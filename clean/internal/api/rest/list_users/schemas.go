package list_users

import (
	"time"

	"github.com/google/uuid"
)

type UserItem struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Surname    *string   `json:"surname,omitempty"`
	TelegramID int64     `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Response struct {
	Items []UserItem `json:"items"`
}
