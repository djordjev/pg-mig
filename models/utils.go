package models

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func BuildConnector(ctx context.Context, str string) (DBConnection, error) {
	conn, err := pgx.Connect(ctx, str)
	if err != nil {
		return nil, err
	}

	return conn, err
}
