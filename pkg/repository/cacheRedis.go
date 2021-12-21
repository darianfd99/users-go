package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/gomodule/redigo/redis"
)

const ExpirationTime = 60

type UserCacheRedisRepository struct {
	listName string
	pool     *redis.Pool
}

func GetRedisConnection(host string, port string, maxIdle int, idleTimeout int) *redis.Pool {
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := redis.Dial("tcp", addr)

	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout * int(time.Second)),
		Dial: func() (redis.Conn, error) {
			return conn, err
		},
	}
}

func NewUserCacheRedisRepository(pool *redis.Pool, listName string) UserCacheRedisRepository {
	return UserCacheRedisRepository{
		listName: listName,
		pool:     pool,
	}
}

func (r UserCacheRedisRepository) GetAll(ctx context.Context) ([]*proto.User, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return []*proto.User{}, err
	}
	defer conn.Close()

	results, err := redis.Strings(conn.Do("lrange", r.listName, 0, -1))
	if err != nil {
		return []*proto.User{}, err
	}

	usersList := make([]*proto.User, 0, len(results))
	for _, result := range results {
		user := proto.User{}
		err := json.Unmarshal([]byte(result), &user)
		if err != nil {
			return []*proto.User{}, err
		}

		usersList = append(usersList, &user)
	}
	return usersList, nil
}

func (r UserCacheRedisRepository) SetList(ctx context.Context, usersList []*proto.User) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	for _, user := range usersList {
		bytes, err := json.Marshal(user)
		if err != nil {
			return err
		}

		_, err = conn.Do("rpush", r.listName, string(bytes))
		if err != nil {
			return err
		}
	}

	_, err = conn.Do("expire", r.listName, ExpirationTime)
	if err != nil {
		return err
	}

	return nil
}
