package subcommands

import (
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"sort"
	"time"
)

// Log structure for log command
type Log struct {
	CommandBase
}

type logGroup struct {
	timestamp int64
	inDB      *int64
	onFS      *filesystem.MigrationFile
}

// Run displays a log of currently present migrations
func (log *Log) Run() error {

	migrations, err := log.getData()
	if err != nil {
		return fmt.Errorf("log command error: unable to fetch migration data %w", err)
	}

	log.printMigrations(migrations)

	return nil
}

func (log *Log) getData() (migrations []logGroup, err error) {
	inDB := make(map[int64]bool)
	onFS := make(map[int64]filesystem.MigrationFile)
	overall := make(map[int64]bool)

	inDBList, err := log.Models.GetMigrationsList()
	if err != nil {
		return
	}

	onFSList, err := log.Filesystem.GetFileTimestamps(time.Time{}, log.Timer.Now())
	if err != nil {
		return
	}

	// Create maps
	for _, migDB := range inDBList {
		inDB[migDB] = true
		overall[migDB] = true
	}

	for _, migFS := range onFSList {
		onFS[migFS.Timestamp] = migFS
		overall[migFS.Timestamp] = true
	}

	// Extract to an ordered array
	timestamps := make([]int64, 0, len(overall))
	for ts, _ := range overall {
		timestamps = append(timestamps, ts)
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })

	// Put data into return array
	migrations = make([]logGroup, 0, len(timestamps))
	for _, ts := range timestamps {
		group := logGroup{timestamp: ts}

		fsMigration, ok := onFS[ts]
		if ok {
			group.onFS = &fsMigration
		} else {
			group.onFS = nil
		}

		_, ok = inDB[ts]
		if ok {
			newTs := ts
			group.inDB = &newTs
		} else {
			group.inDB = nil
		}

		migrations = append(migrations, group)
	}

	return
}

func (log *Log) printMigrations(migrations []logGroup) {
	for _, mig := range migrations {
		var db, fs string

		date := time.Unix(mig.timestamp, 0).Format(time.RFC3339)

		if mig.inDB != nil {
			db = fmt.Sprintf("%d", *mig.inDB)
		}

		if mig.onFS != nil {
			fs = mig.onFS.Up
		}

		log.Printer.PrintMigrations(date, fs, db)
	}
}
