package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djordjev/pg-mig/migrations"
)

// run docker postgres
// docker run --name migration-test-db -p 5434:5432 -e POSTGRES_PASSWORD=testee -e POSTGRES_USER=tester -e POSTGRES_DB=testdb -d postgres:15-alpine

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	migrationsUrl := filepath.Join(wd, "../../examples/db/workspace")

	migrator := migrations.NewRunner(
		"localhost",
		"tester:testee",
		"testdb",
		5434,
		migrationsUrl,
	)

	err = migrator.Run([]string{})

	if err != nil {
		fmt.Println("error", err)
		fmt.Println(migrator.GetPrints())
	}
}
