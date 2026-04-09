package create_notification

import "github.com/google/uuid"

type Request struct {
	UserID  uuid.UUID `json:"user_id"`
	Title   string    `json:"title"`
	Text    string    `json:"text"`
	Channel string    `json:"channel"`
}

type Response struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Channel string    `json:"channel"`
}
