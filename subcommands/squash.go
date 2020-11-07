package subcommands

import (
	"flag"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"time"
)

// Squash structure for squash command
type Squash struct {
	CommandBase
}

// Squash merges migration files into one
func (squash *Squash) Run() error {
	flagSet := flag.NewFlagSet("squash", flag.ExitOnError)

	fromStr := flagSet.String("from", "", "Time of first migration that needs to be squashed")
	toStr := flagSet.String("to", "", "Time of the last migration that needs to be squashed")
	help := flagSet.Bool("help", false, "Prints help for squash command")

	err := flagSet.Parse(squash.Flags)
	if err != nil {
		return fmt.Errorf("squash command error: unable to parse program flags %w", err)
	}

	if help != nil && *help == true {
		flagSet.PrintDefaults()
		return nil
	}

	from, err := squash.Timer.ParseTime(*fromStr)
	if err != nil {
		return err
	}

	to, err := squash.Timer.ParseTime(*toStr)
	if err != nil {
		return err
	}

	migrations, inDB, err := squash.getSquash(from, to)
	if err != nil {
		return err
	}

	last := inDB[len(inDB)-1]

	// Safe to squash migrations
	err = squash.Models.SquashMigrations(from.Unix(), to.Unix(), last)
	if err != nil {
		return err
	}

	err = squash.Filesystem.Squash(migrations)
	if err != nil {
		return err
	}

	return nil
}

func (squash *Squash) getSquash(from time.Time, to time.Time) (migrations filesystem.MigrationFileList, inDB []int64, err error) {
	migrations, err = squash.Filesystem.GetFileTimestamps(from, to)
	if err != nil {
		return
	}

	dbMigs, err := squash.Models.GetMigrationsList()
	fromTS := from.Unix()
	toTS := to.Unix()

	inDB = make([]int64, 0, 0)
	for _, mig := range dbMigs {
		if mig >= fromTS && mig <= toTS {
			inDB = append(inDB, mig)
		}

	}

	if err != nil {
		return
	}

	err = squash.validate(migrations, inDB)

	return
}

func (squash *Squash) validate(list filesystem.MigrationFileList, inDB []int64) error {
	inDBMap := make(map[int64]bool)
	for _, v := range inDB {
		inDBMap[v] = true
	}

	for _, file := range list {
		if file.Up == "" || file.Down == "" {
			return fmt.Errorf("up or down migration missing file %d", file.Timestamp)
		}

		if _, ok := inDBMap[file.Timestamp]; ok {
			delete(inDBMap, file.Timestamp)
		} else {
			return fmt.Errorf("migration file with timestamp %d exists in file list but not in database", file.Timestamp)
		}
	}

	if len(inDBMap) > 0 {
		return fmt.Errorf("some files exist in database but there are not migration files for them")
	}

	return nil
}
