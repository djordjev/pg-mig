package models

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/jackc/pgconn"
)

// DBConnection wrapper interface for interactions with db
type DBConnection interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Close(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// Models interface for interaction with database
type Models interface {
	CreateMetaTable() error
	GetMigrationsList() ([]int64, error)
	Execute(ExecutionContext) error
	SquashMigrations(time.Time, time.Time, int64) error
}

type ExecutionContext struct {
	Sql       string
	IsUp      bool
	Timestamp int64
	Name      string
}
