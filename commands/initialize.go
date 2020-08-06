package commands

import (
	"errors"
	"flag"
	"os"
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
	credentials := flagSet.String("credentials", "", "Credentials for logging in on Postgres instance. In form username@password")

	flagSet.Parse(os.Args[2:])

	if *dbName == "" {
		return errors.New("missing database name")
	}

	if *credentials == "" {
		return errors.New("missing credentials in form username@password")
	}

	config := Config{
		Credentials: *credentials,
		DbName:      *dbName,
		DbURL:       *dbURL,
		Path:        *path,
	}

	err = config.Store()
	if err != nil {
		return err
	}

	return nil
}
