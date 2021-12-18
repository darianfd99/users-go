package repository

import (
	"github.com/darianfd99/users-go/pkg/proto"
)

type UserRepository interface {
	Save(proto.User) (string, error)
	GetAll() ([]*proto.User, error)
	Delete(string) (string, error)
}

type UserCacheRepository interface {
	GetAll() ([]*proto.User, error)
}

//go:generate mockery --case=snake --outpkg=repositorymock --output=repositorymock --name=LocalizationRepository
type Repository struct {
	UserRepository
	UserCacheRepository
}

func NewRepository() *Repository {
	return &Repository{
		UserRepository:      NewUserPostgresRepository(),
		UserCacheRepository: NewUserCacheRedisRepository(),
	}
}
