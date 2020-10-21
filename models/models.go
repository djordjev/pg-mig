package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

const tableName = "__pg_mig_meta"

// ImplModels implementation of models interface with underlying database
type ImplModels struct {
	Db DBConnection
}

// CreateMetaTable creates table named __pg_mig_meta
// that will be used for storing migration info
func (models *ImplModels) CreateMetaTable() error {
	db := models.Db

	_, err := db.Exec(context.Background(), fmt.Sprintf(createMetaTableQuery, tableName))
	if err != nil {
		return fmt.Errorf("db error: unable to create meta table %w", err)
	}

	return nil
}

// GetMigrationsList - fetches timestamps of migrations that has
// been executed in current DB
func (models *ImplModels) GetMigrationsList() ([]int64, error) {
	rows, err := models.Db.Query(context.Background(), fmt.Sprintf(getMigrationsListQuery, tableName))
	if err != nil {
		return nil, fmt.Errorf("db error: unable to query for migrations list %w", err)
	}
	defer rows.Close()

	result := make([]int64, 0, 10)

	for rows.Next() {
		var ts time.Time

		err = rows.Scan(&ts)
		if err != nil {
			return result, fmt.Errorf("db error: unable to scan returned rows from meta table %w", err)
		}

		result = append(result, ts.Unix())
	}

	return result, nil
}

// SquashMigrations deletes all migration instances in meta table between given timestamps (both inclusive).
// and writes a new squash migration with timestamp set to `to` variable value
func (models *ImplModels) SquashMigrations(from int64, to int64) error {
	tx, err := models.Db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("unable to start transaction %w", err)
	}

	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			panic(err)
		}
	}()

	delQuery := fmt.Sprintf("delete from %s where ts >= %d and ts <= %d;", tableName, from, to)

	_, err = tx.Exec(context.Background(), delQuery)
	if err != nil {
		return fmt.Errorf("db error: unable to squash migrations %w", err)
	}

	addQuery := fmt.Sprintf("insert into %s (ts) values ($1);", tableName)
	_, err = tx.Exec(context.Background(), addQuery, time.Unix(to, 0))
	if err != nil {
		return fmt.Errorf("db error: unable to write squash migration %w", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}

func updateMetaTable(executionContext *ExecutionContext, tx pgx.Tx) error {
	unixTs := time.Unix(executionContext.Timestamp, 0)

	upQuery := fmt.Sprintf("insert into %s (ts) values ($1);", tableName)
	downQuery := fmt.Sprintf("delete from %s where ts = $1", tableName)

	var err error
	if executionContext.IsUp {
		_, err = tx.Exec(context.Background(), upQuery, unixTs)
	} else {
		_, err = tx.Exec(context.Background(), downQuery, unixTs)
	}

	return err
}

// Execute runs a migration within a transaction and updates meta table
func (models *ImplModels) Execute(executionContext ExecutionContext) error {
	tx, err := models.Db.Begin(context.Background())
	if err != nil {
		return err
	}

	err = updateMetaTable(&executionContext, tx)
	if err != nil {
		return fmt.Errorf("db error: unable to update meta table %w", err)
	}

	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			panic(err)
		}
	}()

	_, err = tx.Exec(context.Background(), executionContext.Sql)
	if err != nil {
		return fmt.Errorf("db error: unable to execute migration file %s. Error returned %w", executionContext.Name, err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}
