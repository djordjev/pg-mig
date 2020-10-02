package subcommands

import (
	"flag"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"time"
)

// Run structure for run command
type Run struct {
	CommandBase
	GetNow TimeGetter
}

// Run executes up/down migrations
func (run *Run) Run() error {
	flagSet := flag.NewFlagSet("run", flag.ExitOnError)

	strTime := flagSet.String("time", run.GetNow().Format(time.RFC3339), "Time on which you want to upgrade/downgrade DB. Omit for current time")

	// TODO check file formats and matching down files
	inDB, err := run.Models.GetMigrationsList()
	if err != nil {
		return err
	}

	inputTime, err := run.parseTime(strTime)
	if err != nil {
		return err
	}

	stay, down, err := run.getMigrationFiles(inputTime)
	if err != nil {
		return err
	}

	err = run.executeUpMigrations(stay, inDB)
	if err != nil {
		return err
	}

	err = run.executeDownMigrations(down, inDB, inputTime)
	if err != nil {
		return err
	}

	return nil
}

func (run *Run) getMigrationFiles(border time.Time) (stay, goDown filesystem.MigrationFileList, err error) {
	stay, err = run.Filesystem.GetFileTimestamps(time.Time{}, border)
	if err != nil {
		return nil, nil, err
	}

	goDown, err = run.Filesystem.GetFileTimestamps(border, run.GetNow())
	if err != nil {
		return nil, nil, err
	}

	return
}

func (run *Run) parseTime(inputTime *string) (time.Time, error) {
	if inputTime == nil {
		return run.GetNow(), nil
	}

	t, err := time.Parse(time.RFC3339, *inputTime)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (run *Run) executeUpMigrations(stay filesystem.MigrationFileList, inDB []int64) error {
	executedMap := make(map[int64]bool)
	for _, mig := range inDB {
		executedMap[mig] = true
	}

	for _, mig := range stay {
		_, exists := executedMap[mig.Timestamp]
		if exists {
			continue
		}

		migrationContent, err := run.Filesystem.ReadMigrationContent(mig, filesystem.DirectionUp, run.Config)
		if err != nil {
			return err
		}

		err = run.Execute(migrationContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (run *Run) executeDownMigrations(down filesystem.MigrationFileList, inDB []int64, border time.Time) error {
	toExecuteMap := make(map[int64]filesystem.MigrationFile)

	for _, mig := range down {
		toExecuteMap[mig.Timestamp] = mig
	}

	for i := len(inDB) - 1; i >= 0; i-- {
		current := inDB[i]

		if current < border.Unix() {
			// TODO check here should it be < or <=
			return nil
		}

		migrationContent, err := run.Filesystem.ReadMigrationContent(toExecuteMap[current], filesystem.DirectionDown, run.Config)
		if err != nil {
			return err
		}

		err = run.Execute(migrationContent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (run *Run) Execute(content string) error {
	fmt.Println("Executing: ", content)
	return nil
}
