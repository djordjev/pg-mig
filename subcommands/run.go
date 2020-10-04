package subcommands

import (
	"errors"
	"flag"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
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

	strTime := flagSet.String("time", "", "Time on which you want to upgrade/downgrade DB. Omit for current time")

	// TODO check file formats and matching down files
	inDB, err := run.Models.GetMigrationsList()
	if err != nil {
		return err
	}

	inputTime, err := run.parseTime(strTime, inDB)
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

func (run *Run) parseTime(inputTime *string, inDB []int64) (time.Time, error) {
	if inputTime == nil || *inputTime == "" {
		return run.GetNow(), nil
	}

	// Check special values
	if *inputTime == POP {
		if len(inDB) > 0 {
			last := inDB[len(inDB)-1]
			lastTime := time.Unix(last-1, 0) // reduce 1 second from the last migration
			return lastTime, nil
		}
		return time.Time{}, errors.New("pop on empty DB. This is no-op")
	}

	if *inputTime == PUSH {
		var last int64
		if len(inDB) > 0 {
			last = inDB[len(inDB)-1]
		}

		_, down, err := run.getMigrationFiles(time.Unix(last, 0))
		if err != nil {
			return time.Time{}, err
		}
		if len(down) == 0 {
			return time.Time{}, errors.New("push without next migration file. This is no-op")
		}
		return time.Unix(down[0].Timestamp, 0), nil
	}

	// Parse regular timestamp
	t, err := parseAnyTime(*inputTime, run.GetNow)
	if err != nil {
		return t, err
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

		content, err := run.Filesystem.ReadMigrationContent(mig, filesystem.DirectionUp, run.Config)
		if err != nil {
			return err
		}

		execContext := models.ExecutionContext{
			Sql:       content,
			IsUp:      true,
			Timestamp: mig.Timestamp,
			Name:      mig.Up,
		}

		err = run.Models.Execute(execContext)
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

		execContext := models.ExecutionContext{
			Sql:       content,
			IsUp:      false,
			Timestamp: current.Timestamp,
			Name:      current.Down,
		}

		err = run.Models.Execute(execContext)
		if err != nil {
			return err
		}
	}

	return nil
}
