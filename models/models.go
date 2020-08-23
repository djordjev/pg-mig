package models

import (
	"context"
	"fmt"
)

const tableName = "__pg_mig_meta"

// Models interface as a collection for interactions with storage
type Models struct {
	Db DBConnection
}

// CreateMetaTable creates table named __pg_mig_meta
// that will be used for storing migration info
func (models *Models) CreateMetaTable() error {
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
