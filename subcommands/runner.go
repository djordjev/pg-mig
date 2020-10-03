package subcommands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"os"
)

const cmdInit = "init"
const cmdAdd = "add"
const cmdRun = "run"

// Runner structure used for instantiating selected subcommand
type Runner struct {
	Subcommand string
	Flags      []string
	Fs         filesystem.Filesystem
	Connector  DBConnector
	GetNow     TimeGetter
}

// Run runs command selected from args
func (runner *Runner) Run() error {
	if runner.Subcommand == cmdInit {
		err := runner.createInitFile()
		if err != nil {
			return err
		}
	}

	config, err := runner.Fs.LoadConfig()
	if err != nil {
		return err
	}

	connectionString, err := config.GetConnectionString()
	if err != nil {
		return err
	}

	conn, err := runner.Connector(context.Background(), connectionString)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close(context.Background())
		if err != nil {
			panic("Unable to close connection to DB")
		}
	}()

	base := CommandBase{
		Config:     config,
		Models:     &models.ImplModels{Db: conn},
		Flags:      runner.Flags,
		Filesystem: runner.Fs,
	}

	subcommand, err := runner.getSubcommand(&base)
	if err != nil {
		return err
	}

	return subcommand.Run()
}

func (runner *Runner) getSubcommand(base *CommandBase) (Command, error) {
	switch runner.Subcommand {
	case cmdInit:
		{

			initialize := Initialize{CommandBase: *base}
			return &initialize, nil
		}
	case cmdAdd:
		{
			add := Add{CommandBase: *base, GetNow: runner.GetNow}
			return &add, nil
		}

	case cmdRun:
		{
			run := Run{CommandBase: *base, GetNow: runner.GetNow}
			return &run, nil
		}
	}

	return nil, fmt.Errorf("invalid subcommand %s", runner.Subcommand)
}

func (runner *Runner) createInitFile() error {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := flagSet.String("path", wd, "filesystem path where migration definitions are stored. Default: current directory")
	dbURL := flagSet.String("db", "localhost", "url of running PostgreSQL instance. Default localhost")
	dbName := flagSet.String("name", "", "The name of the database against which migrations will run.")
	credentials := flagSet.String("credentials", "", "Credentials for logging in on Postgres instance. In form username:password")
	useSSL := flagSet.String("ssl", "disable", "Whether or not to use ssl. Defaults to disable.")
	port := flagSet.Int("port", 5432, "Port on which PostgreSQL instance is running. Defaults to 5432")

	err = flagSet.Parse(runner.Flags)
	if err != nil {
		return err
	}

	if *dbName == "" {
		return errors.New("missing database name")
	}

	if *credentials == "" {
		return errors.New("missing credentials in form username:password")
	}

	config := filesystem.Config{
		Credentials: *credentials,
		DbName:      *dbName,
		DbURL:       *dbURL,
		Path:        *path,
		SSL:         *useSSL,
		Port:        *port,
	}

	err = runner.Fs.StoreConfig(config)
	if err != nil {
		return err
	}

	return nil
}
