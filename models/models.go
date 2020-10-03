package models

import (
	"context"
	"fmt"
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
			ts timestamp not null
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

func (models *ImplModels) Execute(sql string) error {
	fmt.Println("executing", sql)
	return nil
}
