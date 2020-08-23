package models

import (
	"context"

	"github.com/jackc/pgconn"
)

// DBConnection wrapper interface for interactions with db
type DBConnection interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}
