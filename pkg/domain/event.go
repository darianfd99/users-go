package domain

import (
	"github.com/darianfd99/users-go/kit/event"
)

const UserCreatedEventType event.Type = "events.user.created"

type UserCreatedEvent struct {
	event.BaseEvent
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserCreatedEvent(id, username, email string) UserCreatedEvent {
	return UserCreatedEvent{
		Id:    id,
		Name:  username,
		Email: email,

		BaseEvent: event.NewBaseEvent(id),
	}
}

func (e UserCreatedEvent) Type() event.Type {
	return UserCreatedEventType
}

func (e UserCreatedEvent) UserID() string {
	return e.ID()
}
