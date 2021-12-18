package user

import (
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/google/uuid"
)

func NewUser(username string, email string) *proto.User {
	uuid := uuid.New().String()

	protoUuid := &proto.Uuid{
		Uuid: uuid,
	}
	return &proto.User{
		Uuid:     protoUuid,
		Username: username,
		Email:    email,
	}
}

func NewUsersList(usersList []*proto.User) *proto.UsersList {
	return &proto.UsersList{
		UserList: usersList,
	}
}
