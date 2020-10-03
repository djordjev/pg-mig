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

	createTable := `
		create table if not exists %s (
			id serial primary key,
			ts timestamptz not null
		)
	`

	_, err := db.Exec(context.Background(), fmt.Sprintf(createTable, tableName))
	if err != nil {
		return err
	}

	return nil
}

// GetMigrationsList - fetches timestamps of migrations that has
// been executed in current DB
func (models *ImplModels) GetMigrationsList() ([]int64, error) {

	query := `
		select ts from __pg_mig_meta order by ts asc
	`

	rows, err := models.Db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]int64, 0, 10)

	for rows.Next() {
		var ts time.Time

		err = rows.Scan(&ts)
		if err != nil {
			return result, err
		}

		result = append(result, ts.Unix())
	}

	return result, nil
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

func (models *ImplModels) Execute(executionContext ExecutionContext) error {
	tx, err := models.Db.Begin(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Executing migration %s", executionContext.Name))

	err = updateMetaTable(&executionContext, tx)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback(context.Background())
		if err != nil && err != pgx.ErrTxClosed {
			panic(err)
		}
	}()

	_, err = tx.Exec(context.Background(), executionContext.Sql)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	return nil
}
