package events

import (
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	EventType EventType
	ID        uuid.UUID
	CreatedAt time.Time
}

func NewBaseEvent(eventType EventType) BaseEvent {
	return BaseEvent{
		EventType: eventType,
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}
}

func (e BaseEvent) Type() EventType {
	return e.EventType
}
