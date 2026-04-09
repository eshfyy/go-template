package domain

import (
	"go-template/pkg/optional"
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string

const (
	NotificationChannelTelegram NotificationChannel = "telegram"
)

func isValidChannel(ch NotificationChannel) bool {
	switch ch {
	case NotificationChannelTelegram:
		return true
	default:
		return false
	}
}

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSuccess NotificationStatus = "success"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type Notification struct {
	BaseEntity
	UserID  uuid.UUID
	Title   string
	Text    string
	Channel NotificationChannel
	Status  NotificationStatus
	SentAt  optional.Optional[time.Time]
}

func NewNotification(
	userID uuid.UUID,
	title string,
	text string,
	channel NotificationChannel,
) (Notification, error) {
	fields := make(map[string]string)
	if userID == uuid.Nil {
		fields["user_id"] = "required"
	}
	if title == "" {
		fields["title"] = "required"
	}
	if text == "" {
		fields["text"] = "required"
	}
	if !isValidChannel(channel) {
		fields["channel"] = "unsupported"
	}
	if len(fields) > 0 {
		return Notification{}, &ValidationError{Fields: fields}
	}

	return Notification{
		BaseEntity: NewBaseEntity(),
		UserID:     userID,
		Title:      title,
		Text:       text,
		Channel:    channel,
		Status:     NotificationStatusPending,
		SentAt:     optional.None[time.Time](),
	}, nil
}

func (n *Notification) MarkSent() {
	n.Status = NotificationStatusSuccess
	n.SentAt = optional.Some(time.Now().UTC())
}

func (n *Notification) MarkFailed() {
	n.Status = NotificationStatusFailed
}

const NotificationCreated EventType = "notification_created"

type NotificationCreatedEvent struct {
	BaseEvent
	NotificationID uuid.UUID
}

func (e NotificationCreatedEvent) Key() string {
	return e.NotificationID.String()
}

func NewNotificationCreatedEvent(notificationID uuid.UUID) NotificationCreatedEvent {
	return NotificationCreatedEvent{
		BaseEvent:      NewBaseEvent(NotificationCreated),
		NotificationID: notificationID,
	}
}
