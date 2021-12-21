package main

import (
	"context"
	"log"
	"time"

	"github.com/darianfd99/users-go/pkg/domain"
	"github.com/darianfd99/users-go/pkg/repository"
	"github.com/darianfd99/users-go/pkg/repository/eventBus"
	"github.com/darianfd99/users-go/pkg/server"
	"github.com/darianfd99/users-go/pkg/service"
)

const (
	redisListName         = "users"
	redisHost             = "redis"
	redisPort             = "6379"
	redisMaxIdle          = 3
	redisIdleTimeout      = 400
	psqlTable             = "go-users"
	serverPort            = "8080"
	serverShutdownTimeout = 200
)

func init() {
	time.Sleep(5 * time.Second)
}

func main() {

	//Enviroment
	postgresConfig := repository.PostgresConfig{
		Host:     "postgres",
		Username: "postgres",
		Password: "postgres",
		Port:     "5432",
		Database: "public",
		Sslmode:  "disable",
	}

	rabbitConfig := eventBus.RabbitConfig{
		Host:     "rabbit",
		Username: "guest",
		Password: "guest",
		Port:     "5672",
	}

	db, err := repository.GetPostgresConnection(postgresConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10000)
	db.SetConnMaxLifetime(4 * time.Second)
	db.SetMaxIdleConns(10000)

	pool := repository.GetRedisConnection(redisHost, redisPort, redisMaxIdle, redisIdleTimeout)
	repo := repository.NewRepository(db, psqlTable, pool, redisListName)

	//Rabbit connection for pusher
	conn, ch := eventBus.GetRabbitConn(rabbitConfig)
	defer conn.Close()
	defer ch.Close()

	eventBusP := eventBus.NewEventBusRabbit(ch)

	services := service.NewService(repo, eventBusP)

	ctx, srv := server.NewServer(context.Background(), serverPort, *services, serverShutdownTimeout*time.Second)

	//Rabbit connection for subscriber
	connC, chC := eventBus.GetRabbitConn(rabbitConfig)
	defer connC.Close()
	defer chC.Close()

	eventConsumer := eventBus.NewEventBusRabbit(chC)
	eventConsumer.Subscribe(
		ctx,
		domain.UserCreatedEventType,
		"consumer1")

	srv.Run(ctx)
}
