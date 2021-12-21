package clickhouseLog

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/darianfd99/users-go/pkg/domain"
)

type ClickHouseLogRepository struct {
	connect *sql.DB
}

func GetClickhouseConnection(host string, port string) (*sql.DB, error) {
	addr := fmt.Sprintf("tcp://%s:%s", host, port)
	connect, err := sql.Open("clickhouse", addr)
	if err != nil {
		log.Fatal(err)
	}

	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			return nil, err
		}
		return nil, err
	}

	return connect, nil
}

func NewClickHouseLogRepository(connect *sql.DB) (*ClickHouseLogRepository, error) {

	_, err := connect.Exec(`
		CREATE TABLE IF NOT EXISTS log (
			eventID     String,
			aggregateID String,
			occurredOn  Datetime,
			id    		String,
			name  		String,
			email 		String 
		) engine=MergeTree()
		ORDER BY occurredOn
	`)

	if err != nil {
		return &ClickHouseLogRepository{}, err
	}
	return &ClickHouseLogRepository{connect: connect}, nil
}

func (r *ClickHouseLogRepository) Save(evt domain.UserCreatedEvent) error {
	tx, err := r.connect.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO log (eventID, aggregateID, occurredOn, id, name, email) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	stmt.Exec(evt.ID(), evt.AggregateID(), evt.OccurredOn(), evt.Id, evt.Name, evt.Email)

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
