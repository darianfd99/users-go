package service

import (
	"context"
	"log"
	"sync"

	"github.com/darianfd99/users-go/pkg/domain"
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/darianfd99/users-go/pkg/repository"
)

type UserService struct {
	lock     *sync.Mutex
	repo     repository.UserRepository
	cache    repository.UserCacheRepository
	eventBus repository.EventRepository
}

func NewUserService(repos *repository.Repository, eventBus repository.EventRepository) *UserService {
	return &UserService{
		lock:     &sync.Mutex{},
		repo:     repos.UserRepository,
		cache:    repos.UserCacheRepository,
		eventBus: eventBus,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.RequestUser) (*proto.User, error) {

	user, err := domain.NewUser(req.Username, req.Email)
	if err != nil {
		return &proto.User{}, err
	}

	err = s.repo.Save(ctx, user)
	if err != nil {
		return &proto.User{}, err
	}
	evts := user.PullEvents()

	go func() {
		s.lock.Lock()
		err = s.eventBus.Publish(ctx, evts)
		if err != nil {
			log.Fatal(err)
		}
		s.lock.Unlock()
	}()

	protoUuid := proto.Uuid{
		Uuid: user.GetUuid(),
	}
	protoUser := &proto.User{
		Uuid:     &protoUuid,
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
	}
	return protoUser, nil

}

func (s *UserService) GetAllUsers(ctx context.Context, null *proto.Null) (*proto.UsersList, error) {
	cacheUsersList, err := s.cache.GetAll(ctx)
	if err != nil {
		return &proto.UsersList{}, err
	}

	if len(cacheUsersList) > 0 {
		protoCacheUsersList := domain.NewUsersList(cacheUsersList)
		return protoCacheUsersList, nil
	}

	usersList, err := s.repo.GetAll(ctx)
	if err != nil {
		return &proto.UsersList{}, err
	}

	err = s.cache.SetList(ctx, usersList)
	if err != nil {
		return &proto.UsersList{}, err
	}

	protoUsersList := domain.NewUsersList(usersList)
	return protoUsersList, nil
}

func (s *UserService) EliminateUser(ctx context.Context, uuid *proto.Uuid) (*proto.Uuid, error) {
	err := s.repo.Delete(ctx, uuid.GetUuid())
	if err != nil {
		return &proto.Uuid{}, err
	}

	return uuid, nil
}
