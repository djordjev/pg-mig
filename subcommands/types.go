package subcommands

import (
	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/utils"
	"github.com/spf13/afero"
)

// Command interface encapsulating different commands
type Command interface {
	Run() error
}

// CommandBase base struct for command
type CommandBase struct {
	Models     models.Models
	Config     utils.Config
	Flags      []string
	Filesystem afero.Fs
}