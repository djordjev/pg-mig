package subcommands

import (
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/timer"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type timerArgs struct {
	date string
	fs   string
	db   string
}

func TestLog(t *testing.T) {
	now := "2020-09-20T15:00:00Z"
	loc, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}

	t1, _ := time.Parse(time.RFC3339, "2020-10-20T10:00:00Z")
	t2, _ := time.Parse(time.RFC3339, "2020-10-21T10:00:00Z")
	t3, _ := time.Parse(time.RFC3339, "2020-10-22T10:00:00Z")
	tNow, _ := time.Parse(time.RFC3339, now)

	table := []struct {
		name      string
		inDB      []int64
		inDBErr   error
		onFS      filesystem.MigrationFileList
		onFSErr   error
		timerArgs []timerArgs
	}{
		{
			name: "prints all",
			inDB: []int64{t1.Unix(), t2.Unix(), t3.Unix()},
			onFS: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: t1.Unix(), Up: "t1_up.sql"},
				filesystem.MigrationFile{Timestamp: t2.Unix(), Up: "t2_up.sql"},
				filesystem.MigrationFile{Timestamp: t3.Unix(), Up: "t3_up.sql"},
			},
			timerArgs: []timerArgs{
				{date: t1.In(loc).Format(time.RFC3339), fs: "t1_up.sql", db: fmt.Sprintf("%d", t1.Unix())},
				{date: t2.In(loc).Format(time.RFC3339), fs: "t2_up.sql", db: fmt.Sprintf("%d", t2.Unix())},
				{date: t3.In(loc).Format(time.RFC3339), fs: "t3_up.sql", db: fmt.Sprintf("%d", t3.Unix())},
			},
		},
		{
			name: "missing second in db",
			inDB: []int64{t1.Unix(), t3.Unix()},
			onFS: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: t1.Unix(), Up: "t1_up.sql"},
				filesystem.MigrationFile{Timestamp: t2.Unix(), Up: "t2_up.sql"},
				filesystem.MigrationFile{Timestamp: t3.Unix(), Up: "t3_up.sql"},
			},
			timerArgs: []timerArgs{
				{date: t1.In(loc).Format(time.RFC3339), fs: "t1_up.sql", db: fmt.Sprintf("%d", t1.Unix())},
				{date: t2.In(loc).Format(time.RFC3339), fs: "t2_up.sql", db: ""},
				{date: t3.In(loc).Format(time.RFC3339), fs: "t3_up.sql", db: fmt.Sprintf("%d", t3.Unix())},
			},
		},
		{
			name: "missing second on fs",
			inDB: []int64{t1.Unix(), t2.Unix(), t3.Unix()},
			onFS: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: t1.Unix(), Up: "t1_up.sql"},
				filesystem.MigrationFile{Timestamp: t3.Unix(), Up: "t3_up.sql"},
			},
			timerArgs: []timerArgs{
				{date: t1.In(loc).Format(time.RFC3339), fs: "t1_up.sql", db: fmt.Sprintf("%d", t1.Unix())},
				{date: t2.In(loc).Format(time.RFC3339), fs: "", db: fmt.Sprintf("%d", t2.Unix())},
				{date: t3.In(loc).Format(time.RFC3339), fs: "t3_up.sql", db: fmt.Sprintf("%d", t3.Unix())},
			},
		},
	}

	fmt.Println("args log_test ->", table[0].timerArgs[0].date)

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			now := buildGetNow(now)

			mp := mockedPrinter{}
			models := mockedModels{}
			fs := mockedFilesystem{}

			models.On("GetMigrationsList").Return(test.inDB, test.inDBErr)
			fs.On("GetFileTimestamps", time.Time{}, tNow).Return(test.onFS, test.onFSErr)

			for _, v := range test.timerArgs {
				mp.On("PrintMigrations", v.date, v.fs, v.db).Once()
			}

			log := Log{
				CommandBase: CommandBase{
					Models:     &models,
					Filesystem: &fs,
					Printer:    &mp,
					Timer:      timer.Timer{Now: now},
				},
			}

			err := log.Run()

			if test.onFSErr == nil && test.inDBErr == nil {
				r.NoError(err)
			} else {
				r.Error(err)
			}
		})
	}
}
