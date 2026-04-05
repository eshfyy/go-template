package events

import (
	"github.com/google/uuid"
)

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
