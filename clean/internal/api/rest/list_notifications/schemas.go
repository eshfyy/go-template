package list_notifications

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type NotificationItem struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Channel   string     `json:"channel"`
	Status    string     `json:"status"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Response struct {
	Items []NotificationItem `json:"items"`
}
