package models

import (
	"context"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgconn"
)

// DBConnection wrapper interface for interactions with db
type DBConnection interface {
	Close(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// Models interface for interaction with database
type Models interface {
	CreateMetaTable() error
	GetMigrationsList() ([]int64, error)
	Execute(string) error
}
