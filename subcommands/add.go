package subcommands

import (
	"flag"
	"fmt"
	"strings"
)

// Add structure for init command
type Add struct {
	CommandBase
}

// Run creates two migration files in path directory
func (add *Add) Run() error {
	flagSet := flag.NewFlagSet("add", flag.ExitOnError)

	name := flagSet.String("name", "", "The name of a new revision. It will be used to construct file name. Filenames will be unique even if left blank.")
	help := flagSet.Bool("help", false, "Prints help for add command")

	err := flagSet.Parse(add.Flags)
	if err != nil {
		return fmt.Errorf("add command error: unable to parse program flags %w", err)
	}

	if help != nil && *help == true {
		flagSet.PrintDefaults()
		return nil
	}

	now := add.Timer.Now()
	ms := now.Unix()

	var nameFormatted string
	if *name != "" {
		escapedName := strings.ReplaceAll(*name, " ", "-")
		underlineEscaped := strings.ReplaceAll(escapedName, "_", "-")
		nameFormatted = fmt.Sprintf("_%s", underlineEscaped)
	}

	upName := fmt.Sprintf("mig_%d%s_up.sql", ms, nameFormatted)
	downName := fmt.Sprintf("mig_%d%s_down.sql", ms, nameFormatted)

	err = add.Filesystem.CreateMigrationFile(upName, add.Config.Path)
	if err != nil {
		return err
	}

	err = add.Filesystem.CreateMigrationFile(downName, add.Config.Path)
	if err != nil {
		return err
	}

	return nil
}
