package dbdriver

import (
	"context"
	"fmt"

	"github.com/djordjev/pg-mig/utils"
	"github.com/jackc/pgx/v4"
)

/*
GetDBConnection - establishes connection to DB with data used from config
*/
func GetDBConnection(config utils.Config) (*pgx.Conn, error) {
	connectionString := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=%s", config.Credentials, config.DbURL, config.Port, config.DbName, config.SSL)

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	conn.Exec(context.Background(), "CREATE TABLE test (id SERIAL PRIMARY KEY )")

	return conn, nil
}
