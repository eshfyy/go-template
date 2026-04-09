package domain

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func NewBaseEntity() BaseEntity {
	return BaseEntity{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}
}

type EventType string

type Event interface {
	Key() string
	Type() EventType
}

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
