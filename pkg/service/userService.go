package service

import (
	"context"
	"errors"

	user "github.com/darianfd99/users-go/pkg"
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/darianfd99/users-go/pkg/repository"
)

type UserService struct {
	repo  repository.UserRepository
	cache repository.UserCacheRepository
}

func NewUserService(repos *repository.Repository) *UserService {
	return &UserService{
		repo:  repos.UserRepository,
		cache: repos.UserCacheRepository,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.RequestUser) (*proto.User, error) {
	if req.Username == "" {
		return &proto.User{}, errors.New("username required")
	}

	if req.Email == "" {
		return &proto.User{}, errors.New("email required")
	}

	user := user.NewUser(req.Username, req.Email)

	_, err := s.repo.Save(*user)
	if err != nil {
		return &proto.User{}, err
	}

	return user, nil

}

func (s *UserService) GetAllUsers(ctx context.Context, null *proto.Null) (*proto.UsersList, error) {
	userList, err := s.repo.GetAll()
	if err != nil {
		return &proto.UsersList{}, err
	}

	protoUserList := user.NewUsersList(userList)

	return protoUserList, nil
}

func (s *UserService) EliminateUser(ctx context.Context, uuid *proto.Uuid) (*proto.Uuid, error) {
	StringUuid, err := s.repo.Delete(uuid.GetUuid())
	if err != nil {
		return &proto.Uuid{}, err
	}

	ResponseUuid := &proto.Uuid{
		Uuid: StringUuid,
	}

	return ResponseUuid, nil
}
