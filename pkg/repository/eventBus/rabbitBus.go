package eventBus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/darianfd99/users-go/kit/event"
	"github.com/darianfd99/users-go/pkg/domain"
	"github.com/darianfd99/users-go/pkg/repository/clickhouseLog"
	"github.com/streadway/amqp"
)

const (
	ClickhouseHost = "clickhouse"
	ClickhousePort = "9000"
)

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func GetRabbitConn(rabbitConfig RabbitConfig) (*amqp.Connection, *amqp.Channel) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitConfig.Username, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return conn, ch
}

type EventBusRabbit struct {
	ch *amqp.Channel
}

func NewEventBusRabbit(ch *amqp.Channel) EventBusRabbit {
	return EventBusRabbit{
		ch: ch,
	}
}

func (e EventBusRabbit) Publish(ctx context.Context, evts []event.Event) error {

	errorsChannel := make(chan error)
	for _, evt := range evts {

		go func(evt event.Event) {
			createdUserEvt, ok := evt.(domain.UserCreatedEvent)
			if !ok {
				err := errors.New("unexpected event")
				log.Fatal(err)
			}
			topic := string(evt.Type())
			_, err := e.ch.QueueDeclare(
				topic,
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				errorsChannel <- err
				return
			}

			msg, err := json.Marshal(createdUserEvt)
			if err != nil {
				errorsChannel <- err
				return
			}

			if err != nil {
				errorsChannel <- err
				return
			}

			err = e.ch.Publish(
				"",
				topic,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				},
			)

			errorsChannel <- err
		}(evt)

	}

	for i := 0; i < len(evts); i++ {
		err := <-errorsChannel
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil

}

func (e EventBusRabbit) Subscribe(ctx context.Context, evtType event.Type, consumerName string) {

	conn, err := clickhouseLog.GetClickhouseConnection(ClickhouseHost, ClickhousePort)
	if err != nil {
		log.Fatal(err)
	}

	logRepository, err := clickhouseLog.NewClickHouseLogRepository(conn)
	if err != nil {
		log.Fatal(err)
	}

	topic := string(evtType)
	_, err = e.ch.QueueDeclare(
		topic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := e.ch.Consume(
		topic,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for m := range msgs {
			bytes := m.Body
			var evt domain.UserCreatedEvent
			err = json.Unmarshal(bytes, &evt)

			logRepository.Save(evt)

		}
	}()

}
