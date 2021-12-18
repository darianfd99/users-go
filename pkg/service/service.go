package service

import (
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/darianfd99/users-go/pkg/repository"
)

type Service struct {
	proto.UsersServiceServer
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		UsersServiceServer: NewUserService(repos),
	}
}
