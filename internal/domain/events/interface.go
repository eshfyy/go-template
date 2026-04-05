package events

type EventType string

type Event interface {
	Key() string
	Type() EventType
}
