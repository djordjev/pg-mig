package subcommands

import (
	"errors"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/timer"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestSquashRun(t *testing.T) {
	table := []struct {
		name         string
		flags        []string
		returnError  bool
		files        filesystem.MigrationFileList
		filesError   error
		inDB         []int64
		inDBError    error
		squashedName int64
		skipSquash   bool
	}{
		{
			name:  "squashes",
			flags: []string{"-from=2020-09-20T15:05:05+07:00", "-to=2020-09-28T15:05:05+07:00"},
			files: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: 1600725600, Up: "1600725600_up", Down: "1600725600_down"},
				filesystem.MigrationFile{Timestamp: 1600812000, Up: "1600812000_up", Down: "1600812000_down"},
				filesystem.MigrationFile{Timestamp: 1600898400, Up: "1600898400_up", Down: "1600898400_down"},
			},
			inDB:         []int64{1600725600, 1600812000, 1600898400, 1603490400, 1603490600},
			squashedName: 1600898400,
		},
		{
			name:  "returns error db",
			flags: []string{"-from=2020-09-20T15:05:05+07:00", "-to=2020-09-28T15:05:05+07:00"},
			files: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: 1600725600, Up: "1600725600_up", Down: "1600725600_down"},
				filesystem.MigrationFile{Timestamp: 1600812000, Up: "1600812000_up", Down: "1600812000_down"},
				filesystem.MigrationFile{Timestamp: 1600898400, Up: "1600898400_up", Down: "1600898400_down"},
			},
			inDB:         []int64{1600725600, 1600812000, 1600898400, 1603490400, 1603490600},
			inDBError:    errors.New("err"),
			squashedName: 1600898400,
			returnError:  true,
		},
		{
			name:  "returns error fs",
			flags: []string{"-from=2020-09-20T15:05:05+07:00", "-to=2020-09-28T15:05:05+07:00"},
			files: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: 1600725600, Up: "1600725600_up", Down: "1600725600_down"},
				filesystem.MigrationFile{Timestamp: 1600812000, Up: "1600812000_up", Down: "1600812000_down"},
				filesystem.MigrationFile{Timestamp: 1600898400, Up: "1600898400_up", Down: "1600898400_down"},
			},
			inDB:         []int64{1600725600, 1600812000, 1600898400, 1603490400, 1603490600},
			filesError:   errors.New("err"),
			squashedName: 1600898400,
			returnError:  true,
		},
		{
			name:  "migrations mismatch",
			flags: []string{"-from=2020-09-20T15:05:05+07:00", "-to=2020-09-28T15:05:05+07:00"},
			files: filesystem.MigrationFileList{
				filesystem.MigrationFile{Timestamp: 1600725600, Up: "1600725600_up", Down: "1600725600_down"},
				filesystem.MigrationFile{Timestamp: 1600812000, Up: "1600812000_up", Down: "1600812000_down"},
				filesystem.MigrationFile{Timestamp: 1600898400, Up: "1600898400_up", Down: "1600898400_down"},
			},
			inDB:        []int64{1600725600, 1600898400, 1603490400, 1603490600},
			returnError: true,
			skipSquash:  true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)
			mockedFS := mockedFilesystem{}
			mockedMod := mockedModels{}

			squash := Squash{
				CommandBase: CommandBase{
					Filesystem: &mockedFS,
					Flags:      test.flags,
					Models:     &mockedMod,
					Timer:      timer.Timer{Now: buildGetNow("2020-09-20T15:00:00Z")},
				},
			}

			fromStr := strings.Split(test.flags[0], "=")
			from, _ := time.Parse(time.RFC3339, fromStr[1])

			toStr := strings.Split(test.flags[1], "=")
			to, _ := time.Parse(time.RFC3339, toStr[1])

			mockedFS.On("GetFileTimestamps", from.Add(-1*time.Millisecond), to).
				Return(test.files, test.filesError).Once()

			mockedMod.On("GetMigrationsList").Return(test.inDB, test.inDBError).Once()

			if !test.skipSquash {
				mockedFS.On("Squash", test.files).Return(nil).Once()
				mockedMod.On("SquashMigrations", from, to, test.squashedName).
					Return(nil).Once()
			}

			err := squash.Run()
			if test.returnError {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
