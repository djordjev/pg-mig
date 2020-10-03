package subcommands

import (
	"flag"
	"fmt"
)

// Add structure for init command
type Add struct {
	CommandBase
	GetNow TimeGetter
}

// Run creates two migration files in path directory
func (add *Add) Run() error {
	flagSet := flag.NewFlagSet("add", flag.ExitOnError)

	name := flagSet.String("name", "", "The name of a new revision. It will be used to construct file name. Filenames will be unique even if left blank.")

	err := flagSet.Parse(add.Flags)
	if err != nil {
		return err
	}

	now := add.GetNow()
	ms := now.Unix()

	var nameFormatted string
	if *name != "" {
		nameFormatted = fmt.Sprintf("_%s", *name)
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
