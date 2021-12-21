package domain

import (
	"errors"

	"github.com/darianfd99/users-go/kit/event"
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/google/uuid"
)

type User struct {
	uuid     string
	username string
	email    string

	events []event.Event
}

func NewUser(username string, email string) (User, error) {
	if username == "" {
		return User{}, errors.New("username required")
	}

	if email == "" {
		return User{}, errors.New("email required")
	}

	uuid := uuid.New().String()

	user := User{
		uuid:     uuid,
		username: username,
		email:    email,
	}

	user.Record(NewUserCreatedEvent(user.uuid, user.username, user.email))

	return user, nil
}

func (u User) GetUuid() string {
	return u.uuid
}

func (u User) GetUsername() string {
	return u.username
}

func (u User) GetEmail() string {
	return u.email
}

// Record records a new domain event.
func (u *User) Record(evt event.Event) {
	u.events = append(u.events, evt)
}

func (u *User) PullEvents() []event.Event {
	evts := u.events
	u.events = []event.Event{}
	return evts
}

func NewUsersList(usersList []*proto.User) *proto.UsersList {
	return &proto.UsersList{
		UserList: usersList,
	}
}
