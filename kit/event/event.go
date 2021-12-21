package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Type string

type Handler interface {
	Handle(context.Context, chan error, Event) error
}

// Event represents a domain command
type Event interface {
	ID() string
	AggregateID() string
	OccurredOn() time.Time
	Type() Type
}

type BaseEvent struct {
	EventID      string    `json:"eventID"`
	EaggregateID string    `json:"aggregateID"`
	EoccurredOn  time.Time `json:"occurredOn"`
}

func NewBaseEvent(aggregateID string) BaseEvent {
	return BaseEvent{
		EventID:      uuid.New().String(),
		EaggregateID: aggregateID,
		EoccurredOn:  time.Now(),
	}
}

func (b BaseEvent) ID() string {
	return b.EventID
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.EoccurredOn
}

func (b BaseEvent) AggregateID() string {
	return b.EaggregateID
}
