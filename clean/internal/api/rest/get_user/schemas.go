package get_user

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Surname    *string   `json:"surname,omitempty"`
	TelegramID int64     `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}
