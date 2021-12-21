package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/darianfd99/users-go/pkg/domain"
	"github.com/darianfd99/users-go/pkg/proto"
	_ "github.com/lib/pq"
)

type UserPostgresRepository struct {
	table string
	db    *sql.DB
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Sslmode  string
}

func GetPostgresConnection(config PostgresConfig) (*sql.DB, error) {
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=%s", config.Username, config.Password, config.Host, config.Port, config.Sslmode)
	return sql.Open("postgres", uri)
}

func NewUserPostgresRepository(db *sql.DB, table string) UserPostgresRepository {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s"(uuid VARCHAR(500) primary key, username VARCHAR(500) unique, email VARCHAR(500) unique);`, table)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return UserPostgresRepository{
		table: table,
		db:    db,
	}
}

func (r UserPostgresRepository) Save(ctx context.Context, user domain.User) error {

	query := fmt.Sprintf(`INSERT INTO "%s"("uuid","username","email") VALUES ($1, $2, $3)`, r.table)
	_, err := r.db.ExecContext(ctx, query, user.GetUuid(), user.GetUsername(), user.GetEmail())
	if err != nil {
		return err
	}

	return nil
}

func (r UserPostgresRepository) GetAll(ctx context.Context) ([]*proto.User, error) {
	query := fmt.Sprintf(`SELECT * FROM "%s"`, r.table)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []*proto.User{}, err
	}
	defer rows.Close()
	usersList := []*proto.User{}
	for rows.Next() {
		user := &proto.User{
			Uuid: &proto.Uuid{},
		}

		err = rows.Scan(&user.Uuid.Uuid, &user.Username, &user.Email)
		if err != nil {
			return []*proto.User{}, err
		}

		usersList = append(usersList, user)
	}
	return usersList, nil
}

func (r UserPostgresRepository) Delete(ctx context.Context, uuid string) error {

	query := `DELETE FROM "$1" WHERE uuid=$2`
	_, err := r.db.ExecContext(ctx, query, r.table, uuid)
	if err != nil {
		return err
	}

	return nil
}
