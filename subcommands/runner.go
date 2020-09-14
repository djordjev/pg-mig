package subcommands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/utils"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
)

const cmdInit = "init"

// Runner structure used for instantiating selected subcommand
type Runner struct {
	Subcommand string
	Flags      []string
	Filesystem afero.Fs
}

// Run runs command selected from args
func (runner *Runner) Run() error {
	if runner.Subcommand == cmdInit {
		err := runner.createInitFile()
		if err != nil {
			return err
		}
	}

	config := utils.Config{Filesystem: runner.Filesystem}
	err := config.Load()
	if err != nil {
		return err
	}

	connectionString, err := config.GetConnectionString()
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	base := CommandBase{
		Filesystem: runner.Filesystem,
		Config:     config,
		Models:     models.Models{Db: conn},
		Flags:      runner.Flags,
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

			initialize := Initialize{
				CommandBase: *base,
			}

			return &initialize, nil
		}
	}

	return nil, fmt.Errorf("Invalid subcommand %s", runner.Subcommand)
}

func (runner *Runner) createInitFile() error {
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

	flagSet.Parse(runner.Flags)

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

		Filesystem: runner.Filesystem,
	}

	err = config.Store()
	if err != nil {
		return err
	}

	return nil
}
