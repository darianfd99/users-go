package service

import (
	"context"

	"github.com/darianfd99/users-go/kit/event"
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/darianfd99/users-go/pkg/repository"
)

type Event interface {
	Publish(ctx context.Context, evt event.Event) error
	Subscribe(ctx context.Context) error
}

type Service struct {
	proto.UsersServiceServer
}

func NewService(repos *repository.Repository, eventBus repository.EventRepository) *Service {
	return &Service{
		UsersServiceServer: NewUserService(repos, eventBus),
	}
}
