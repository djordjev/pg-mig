package migrations

import (
	"context"
	"time"

	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/subcommands"
	"github.com/djordjev/pg-mig/timer"
	"github.com/spf13/afero"
)

var ConfigDirEnv = ""

type migrations struct {
	fs      filesystem.Filesystem
	printer *bufferedPrinter
	config  filesystem.Config
}

// Creates new migration runner to control execution running
func NewRunner(host string, credentials string, dbName string, port int, workDir string) migrations {

	config := filesystem.Config{
		Credentials: credentials,
		DbName:      dbName,
		DbURL:       host,
		Path:        workDir,
		SSL:         "",
		Port:        port,
		NoColor:     true,
	}

	fs := &filesystem.ImplFilesystem{
		Fs:             afero.NewOsFs(),
		GetNow:         time.Now,
		ConfigDir:      workDir,
		ExternalConfig: &config,
	}

	return migrations{fs: fs, printer: newBufferedPrinter(), config: config}
}

func (m migrations) GetPrints() string {
	return m.printer.GetAllPrints()
}

func (m migrations) Run(params []string) error {

	connectionString, err := m.config.GetConnectionString()
	if err != nil {
		return err
	}

	conn, err := models.BuildConnector(context.Background(), connectionString)
	if err != nil {
		return err
	}

	defer func() {
		err := conn.Close(context.Background())
		if err != nil {
			panic("Unable to close connection to DB")
		}
	}()

	base := subcommands.CommandBase{
		Config:     m.config,
		Models:     &models.ImplModels{Db: conn},
		Flags:      params,
		Filesystem: m.fs,
		Timer:      timer.Timer{Now: time.Now},
		Printer:    m.printer,
	}

	init := subcommands.Initialize{CommandBase: base}
	err = init.Run()
	if err != nil {
		return err
	}

	runner := subcommands.Run{CommandBase: base}

	err = runner.Run()
	if err != nil {
		return err
	}

	return nil
}
