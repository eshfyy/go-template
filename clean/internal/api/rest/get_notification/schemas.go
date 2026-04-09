package get_notification

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	Text      string     `json:"text"`
	Channel   string     `json:"channel"`
	Status    string     `json:"status"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
