package queries

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetConnection() (*pgx.Conn) {
	conn, err := pgx.Connect(
		context.Background(),
		os.Getenv("DB_CONNECTION_STRING"),
	)

	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	return conn
}