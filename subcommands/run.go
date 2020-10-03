package subcommands

import (
	"flag"
	"github.com/djordjev/pg-mig/filesystem"
	"sort"
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

	downMigrations := run.getInDBDownMigrations(inDB, inputTime)

	err = run.executeDownMigrations(down, downMigrations)
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

func (run *Run) getInDBDownMigrations(inDB []int64, border time.Time) []int64 {
	result := make([]int64, 0, 10)

	for _, mig := range inDB {
		if mig > border.Unix() {
			result = append(result, mig)
		}
	}

	// Sort in reverse order
	sort.Slice(result, func(i, j int) bool { return result[i] > result[j] })

	return result
}

func (run *Run) parseTime(inputTime *string) (time.Time, error) {
	if inputTime == nil || *inputTime == "" {
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

		err = run.Models.Execute(migrationContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (run *Run) executeDownMigrations(down filesystem.MigrationFileList, downIDs []int64) error {
	if downIDs == nil {
		return nil
	}

	toExecuteMap := make(map[int64]filesystem.MigrationFile)
	config := run.Config

	for _, mig := range down {
		toExecuteMap[mig.Timestamp] = mig
	}

	for _, toExec := range downIDs {
		current := toExecuteMap[toExec]

		content, err := run.Filesystem.ReadMigrationContent(current, filesystem.DirectionDown, config)
		if err != nil {
			return err
		}

		err = run.Models.Execute(content)
		if err != nil {
			return err
		}
	}

	return nil
}
