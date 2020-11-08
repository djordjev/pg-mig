package subcommands

import (
	"context"
	"flag"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/timer"
	"os"
)

const cmdInit = "init"
const cmdAdd = "add"
const cmdRun = "run"
const cmdSquash = "squash"
const cmdLog = "log"
const cmdHelp = "help"

// Runner structure used for instantiating selected subcommand
type Runner struct {
	Subcommand string
	Flags      []string
	Fs         filesystem.Filesystem
	Connector  DBConnector
	Timer      timer.Timer
	Printer    Printer
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

	runner.Printer.SetNoColor(config.NoColor)

	connectionString, err := config.GetConnectionString()
	if err != nil {
		return err
	}

	conn, err := runner.Connector(context.Background(), connectionString)
	if err != nil {
		return fmt.Errorf("run error: unable to connect on database using connection string %s", connectionString)
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
		Timer:      runner.Timer,
		Printer:    runner.Printer,
	}

	subcommand, err := runner.getSubcommand(&base)
	if err != nil {
		return err
	}

	err = subcommand.Run()
	if err != nil {
		// Intercept error and print it here
		runner.Printer.PrintError(fmt.Sprintf("%v", err))
		os.Exit(1)
		return nil
	}

	return err
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
			add := Add{CommandBase: *base}
			return &add, nil
		}

	case cmdRun:
		{
			run := Run{CommandBase: *base}
			return &run, nil
		}

	case cmdSquash:
		{
			squash := Squash{CommandBase: *base}
			return &squash, nil
		}
	case cmdLog:
		{
			log := Log{CommandBase: *base}
			return &log, nil
		}
	case cmdHelp:
		{
			help := Help{}
			return &help, nil
		}
	}

	return nil, fmt.Errorf("run error: invalid subcommand %s", runner.Subcommand)
}

func (runner *Runner) createInitFile() error {
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("run error: unable to react current working directory %w", err)
	}

	path := flagSet.String("path", wd, "filesystem path where migration definitions are stored. Default: current directory")
	dbURL := flagSet.String("db", "localhost", "url of running PostgreSQL instance. Default localhost")
	dbName := flagSet.String("name", "", "The name of the database against which migrations will run.")
	credentials := flagSet.String("credentials", "", "Credentials for logging in on Postgres instance. In form username:password")
	useSSL := flagSet.String("ssl", "disable", "Whether or not to use ssl. Defaults to disable.")
	port := flagSet.Int("port", 5432, "Port on which PostgreSQL instance is running. Defaults to 5432")
	noColor := flagSet.Bool("nocolor", false, "prevent pg-mig for printing emojis and colored text. Useful on terminals not supporting unicode.")
	help := flagSet.Bool("help", false, "Prints help for init command")

	err = flagSet.Parse(runner.Flags)
	if err != nil {
		return fmt.Errorf("run error: unable to parse flags %+q", runner.Flags)
	}

	if help != nil && *help == true {
		flagSet.PrintDefaults()
		return nil
	}

	if *dbName == "" {
		return fmt.Errorf("run error: missing database name")
	}

	if *credentials == "" {
		return fmt.Errorf("run error: missing credentials in form username:password")
	}

	config := filesystem.Config{
		Credentials: *credentials,
		DbName:      *dbName,
		DbURL:       *dbURL,
		Path:        *path,
		SSL:         *useSSL,
		Port:        *port,
		NoColor:     *noColor,
	}

	err = runner.Fs.StoreConfig(config)
	if err != nil {
		return err
	}

	return nil
}
