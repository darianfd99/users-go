package repository

import (
	"context"
	"database/sql"

	"github.com/darianfd99/users-go/kit/event"
	"github.com/darianfd99/users-go/pkg/domain"
	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/gomodule/redigo/redis"
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) error
	GetAll(ctx context.Context) ([]*proto.User, error)
	Delete(ctx context.Context, uuid string) error
}

type UserCacheRepository interface {
	GetAll(ctx context.Context) ([]*proto.User, error)
	SetList(ctx context.Context, UserList []*proto.User) error
}

type EventRepository interface {
	Publish(ctx context.Context, evts []event.Event) error
	Subscribe(ctx context.Context, evtType event.Type, consumerName string)
}

//go:generate mockery --case=snake --outpkg=repositorymock --output=repositorymock --name=LocalizationRepository
type Repository struct {
	UserRepository
	UserCacheRepository
}

func NewRepository(db *sql.DB, table string, pool *redis.Pool, listName string) *Repository {
	return &Repository{
		UserRepository:      NewUserPostgresRepository(db, table),
		UserCacheRepository: NewUserCacheRedisRepository(pool, listName),
	}
}
