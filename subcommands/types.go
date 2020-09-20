package subcommands

import (
	"context"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/spf13/afero"
)

// Command interface encapsulating different commands
type Command interface {
	Run() error
}

// CommandBase base struct for command
type CommandBase struct {
	Models     models.Models
	Config     filesystem.Config
	Flags      []string
	Filesystem afero.Fs
}

// DBConnector interface for opening DB connection
type DBConnector func(ctx context.Context, connString string) (models.DBConnection, error)
