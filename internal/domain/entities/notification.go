package entities

import (
	vo "go-template/internal/domain/valueobjects"
	"go-template/pkg/optional"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	BaseEntity
	UserID  uuid.UUID
	Title   string
	Text    string
	Channel vo.NotificationChannel
	Status  vo.NotificationStatus
	SentAt  optional.Optional[time.Time]
}

func NewNotification(
	userID uuid.UUID,
	title string,
	text string,
	channel vo.NotificationChannel,
) Notification {
	return Notification{
		BaseEntity: NewBaseEntity(),
		UserID:     userID,
		Title:      title,
		Text:       text,
		Channel:    channel,
		Status:     vo.NotificationStatusPending,
		SentAt:     optional.None[time.Time](),
	}
}

func (n *Notification) SendNotification(
	status vo.NotificationStatus,
	SentAt optional.Optional[time.Time],
) {
	n.SentAt = SentAt
	n.Status = status
}
