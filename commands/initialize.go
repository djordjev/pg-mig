package commands

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/djordjev/pg-mig/dbdriver"
	"github.com/djordjev/pg-mig/utils"
)

/*
Initialize function creates config file
*/
func Initialize() error {

	flagSet := flag.NewFlagSet("init", flag.ExitOnError)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := flagSet.String("path", wd, "filesystem path where migration definitions are stored. Default: current directory")
	dbURL := flagSet.String("db", "localhost", "url of running PostgreSQL instance. Default localhost")
	dbName := flagSet.String("name", "", "The name of the database agains which migrations will run.")
	credentials := flagSet.String("credentials", "", "Credentials for logging in on Postgres instance. In form username:password")
	useSSL := flagSet.String("ssl", "disable", "Wheather or not to use ssl. Defaults to disable.")
	port := flagSet.Int("port", 5432, "Port on which PostgreSQL instance is running. Defaults to 5432")

	flagSet.Parse(os.Args[2:])

	if *dbName == "" {
		return errors.New("missing database name")
	}

	if *credentials == "" {
		return errors.New("missing credentials in form username:password")
	}

	config := utils.Config{
		Credentials: *credentials,
		DbName:      *dbName,
		DbURL:       *dbURL,
		Path:        *path,
		SSL:         *useSSL,
		Port:        *port,
	}

	err = config.Store()
	if err != nil {
		return err
	}

	db, err := dbdriver.GetDBConnection(config)
	if err != nil {
		return err
	}

	defer db.Close(context.Background())

	return nil
}
