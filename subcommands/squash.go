package subcommands

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
)

// Squash structure for squash command
type Squash struct {
	CommandBase
}

// Squash merges migration files into one
func (squash *Squash) Run() error {
	flagSet := flag.NewFlagSet("squash", flag.ExitOnError)

	from := flagSet.String("from", "", "Time of first migration that needs to be squashed")
	to := flagSet.String("to", "", "Time of the last migration that needs to be squashed")

	err := flagSet.Parse(squash.Flags)
	if err != nil {
		return fmt.Errorf("squash command error: unable to parse program flags %w", err)
	}

	files, err := squash.getMigrationsBetweenTimes(*from, *to)
	if err != nil {
		return err
	}

	//err = squash.Models.SquashMigrations(from, to)
	//if err != nil {
	//	return err
	//}
	//
	//up, err := squash.mergeUpFiles(files)
	//if err != nil {
	//	return err
	//}
	//
	//down, err := squash.mergeDownFiles(files)
	//if err != nil {
	//	return err
	//}
	//
	//squash.Filesystem.DeleteMigrationFiles(files)
	//squash.Filesystem.CreateMigrationFile()
	fmt.Println(files)
	// TODO merge files
	return nil
}

func (squash *Squash) getMigrationsBetweenTimes(fromTime string, toTime string) (filesystem.MigrationFileList, error) {
	from, err := squash.Timer.ParseTime(fromTime)
	if err != nil {
		return nil, err
	}

	to, err := squash.Timer.ParseTime(toTime)
	if err != nil {
		return nil, err
	}

	inDB, err := squash.Models.GetMigrationsList()
	if err != nil {
		return nil, err
	}

	files, err := squash.Filesystem.GetFileTimestamps(from, to)
	if err != nil {
		return nil, err
	}

	err = squash.validate(files, inDB)
	if err != nil {
		return nil, err
	}

	return files, nil
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

func (squash *Squash) mergeUpFiles(list filesystem.MigrationFileList) (string, error) {
	result := bytes.NewBufferString("")

	for _, file := range list {
		content, err := squash.Filesystem.ReadMigrationContent(file, filesystem.DirectionUp, squash.Config)
		if err != nil {
			return "", err
		}

		result.WriteString(fmt.Sprintf("-- migration %d UP\n", file.Timestamp))
		result.WriteString(content)
		result.WriteString("\n")
	}

	return result.String(), nil
}

func (squash *Squash) mergeDownFiles(list filesystem.MigrationFileList) (string, error) {
	result := bytes.NewBufferString("")

	for i := len(list) - 1; i >= 0; i-- {
		file := list[i]

		content, err := squash.Filesystem.ReadMigrationContent(file, filesystem.DirectionDown, squash.Config)
		if err != nil {
			return "", err
		}

		result.WriteString(fmt.Sprintf("-- migration %d DOWN\n", file.Timestamp))
		result.WriteString(content)
		result.WriteString("\n")
	}

	return result.String(), nil
}
