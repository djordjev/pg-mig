package models

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

const reconnectCount = 3
const reconnectTimeout = 3

func BuildConnector(ctx context.Context, str string) (conn DBConnection, err error) {
	for i := 0; i < reconnectCount; i++ {
		conn, err := pgx.Connect(ctx, str)
		if err == nil {
			return conn, err
		}

		time.Sleep(reconnectTimeout * time.Second)
	}

	return
}
