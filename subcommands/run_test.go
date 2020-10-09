package subcommands

import (
	"errors"
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/timer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getNow() time.Time {
	now, _ := time.Parse(time.RFC3339, "2020-09-20T15:00:00Z")
	return now
}

func TestGetMigrationFiles(t *testing.T) {
	r := require.New(t)

	upList := filesystem.MigrationFileList{
		filesystem.MigrationFile{Timestamp: 1, Up: "1_up", Down: "1_down"},
		filesystem.MigrationFile{Timestamp: 2, Up: "2_up", Down: "2_down"},
		filesystem.MigrationFile{Timestamp: 3, Up: "3_up", Down: "3_down"},
	}

	downList := filesystem.MigrationFileList{
		filesystem.MigrationFile{Timestamp: 4, Up: "4_up", Down: "4_down"},
		filesystem.MigrationFile{Timestamp: 5, Up: "5_up", Down: "5_down"},
		filesystem.MigrationFile{Timestamp: 6, Up: "6_up", Down: "6_down"},
	}

	mockFS := &mockedFilesystem{}

	run := Run{
		CommandBase: CommandBase{
			Models:     mockedModels{},
			Filesystem: mockFS,
			Timer:      timer.Timer{Now: getNow},
		},
	}

	now := getNow()
	border := time.Unix(123, 0)
	mockFS.On("GetFileTimestamps", time.Time{}, border).Return(upList, nil)
	mockFS.On("GetFileTimestamps", border, now).Return(downList, nil)

	stay, down, err := run.getMigrationFiles(border)
	for k, v := range stay {
		r.Equal(v, upList[k], "up migrations are not matching")
	}

	for k, v := range down {
		r.Equal(v, downList[k], "down migrations are not matching")
	}

	mockFS.AssertExpectations(t)

	r.Nil(err, "returned error")

	newBorder := time.Unix(321, 0)
	mockFS.On("GetFileTimestamps", time.Time{}, newBorder).Return(filesystem.MigrationFileList{}, errors.New("err"))
	_, _, err = run.getMigrationFiles(newBorder)

	r.NotNil(err, "should return error")
}

func TestParseTime(t *testing.T) {
	r := require.New(t)

	t1, _ := time.Parse(time.RFC3339, "2020-09-20T13:00:00Z")

	table := []struct {
		val      string
		expected time.Time
		inDB     []int64
		err      bool
		mockMF   filesystem.MigrationFileList
	}{
		{
			val:      "",
			expected: getNow(),
			inDB:     []int64{},
		},
		{
			val:      "2020-09-20T13:00:00Z",
			expected: t1,
			inDB:     []int64{},
		},
		{
			val:      "invalid date",
			expected: time.Time{},
			inDB:     []int64{},
			err:      true,
		},
		{
			val:      "pop",
			expected: time.Unix(2, 0),
			inDB: []int64{
				time.Unix(1, 0).Unix(),
				time.Unix(2, 0).Unix(),
				time.Unix(3, 0).Unix(),
			},
		},
		{
			val:      "push",
			expected: time.Unix(11, 0),
			mockMF:   filesystem.MigrationFileList{filesystem.MigrationFile{Timestamp: 11}},
			inDB: []int64{
				time.Unix(1, 0).Unix(),
				time.Unix(2, 0).Unix(),
				time.Unix(3, 0).Unix(),
			},
		},
	}

	for _, v := range table {
		t.Run(fmt.Sprintf("test time = %s", v.val), func(t *testing.T) {
			mockedFS := mockedFilesystem{}
			run := Run{
				CommandBase: CommandBase{
					Filesystem: &mockedFS,
					Timer:      timer.Timer{Now: getNow},
				},
			}

			if v.mockMF != nil {
				last := v.inDB[len(v.inDB)-1]
				lastTime := time.Unix(last, 0)
				emptyMFL := filesystem.MigrationFileList{}
				mockedFS.On("GetFileTimestamps", time.Time{}, lastTime).Return(emptyMFL, nil)
				mockedFS.On("GetFileTimestamps", lastTime, getNow()).Return(v.mockMF, nil)
			}

			res, err := run.parseTime(&v.val, v.inDB)
			r.Equal(res, v.expected, "return value mismatch")
			if v.err {
				r.NotNil(err, "error mismatch")
			}
		})
	}
}

func TestExecuteUpMigrations(t *testing.T) {
	mf1 := filesystem.MigrationFile{Timestamp: 1, Up: "1_up.sql", Down: "1_down.sql"}
	mf2 := filesystem.MigrationFile{Timestamp: 2, Up: "2_up.sql", Down: "2_down.sql"}
	mf3 := filesystem.MigrationFile{Timestamp: 3, Up: "3_up.sql", Down: "3_down.sql"}
	mf4 := filesystem.MigrationFile{Timestamp: 4, Up: "4_up.sql", Down: "4_down.sql"}

	table := []struct {
		name              string
		stay              filesystem.MigrationFileList
		inDB              []int64
		shouldReturnError bool
		expectedToRun     []filesystem.MigrationFile
	}{
		{
			name:              "execute migrations when db is empty",
			stay:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			inDB:              []int64{},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{mf1, mf2, mf3, mf4},
		},
		{
			name:              "no up migrations to execute",
			stay:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			inDB:              []int64{1, 2, 3, 4},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{},
		},
		{
			name:              "run one migration",
			stay:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			inDB:              []int64{1, 2, 4},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{mf3},
		},
		{
			name:              "run two migration",
			stay:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			inDB:              []int64{2, 4},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{mf1, mf3},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			r := require.New(t)

			fs := &mockedFilesystem{}
			m := &mockedModels{}

			run := Run{
				CommandBase: CommandBase{
					Filesystem: fs,
					Models:     m,
				},
			}

			// Setup expectations
			for _, migration := range v.expectedToRun {
				fs.On("ReadMigrationContent", migration, mock.Anything, run.Config).Return(migration.Up, nil)
				expectedExec := models.ExecutionContext{
					Sql:       migration.Up,
					IsUp:      true,
					Timestamp: migration.Timestamp,
					Name:      migration.Up,
				}
				m.On("Execute", expectedExec).Return(nil)
			}

			err := run.executeUpMigrations(v.stay, v.inDB)

			fs.AssertExpectations(t)
			m.AssertExpectations(t)

			if v.shouldReturnError {
				r.NotNil(err, "should return error")
			} else {
				r.Nil(err, "should not return error")
			}
		})
	}
}

func TestExecuteDownMigrations(t *testing.T) {
	mf1 := filesystem.MigrationFile{Timestamp: 1, Up: "1_up.sql", Down: "1_down.sql"}
	mf2 := filesystem.MigrationFile{Timestamp: 2, Up: "2_up.sql", Down: "2_down.sql"}
	mf3 := filesystem.MigrationFile{Timestamp: 3, Up: "3_up.sql", Down: "3_down.sql"}
	mf4 := filesystem.MigrationFile{Timestamp: 4, Up: "4_up.sql", Down: "4_down.sql"}

	table := []struct {
		name              string
		down              filesystem.MigrationFileList
		downIDs           []int64
		shouldReturnError bool
		expectedToRun     []filesystem.MigrationFile
	}{
		{
			name:              "execute migrations when db is empty",
			down:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			downIDs:           []int64{},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{},
		},
		{
			name:              "down all migrations",
			down:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			downIDs:           []int64{1, 2, 3, 4},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{mf1, mf2, mf3, mf4},
		},
		{
			name:              "down two migrations",
			down:              filesystem.MigrationFileList{mf1, mf2, mf3, mf4},
			downIDs:           []int64{3, 4},
			shouldReturnError: false,
			expectedToRun:     []filesystem.MigrationFile{mf3, mf4},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			r := require.New(t)

			fs := &mockedFilesystem{}
			m := &mockedModels{}

			run := Run{
				CommandBase: CommandBase{
					Filesystem: fs,
					Models:     m,
				},
			}

			for _, migration := range v.expectedToRun {
				fs.On("ReadMigrationContent", migration, mock.Anything, run.Config).Return(migration.Down, nil)
				expectedExec := models.ExecutionContext{
					Sql:       migration.Down,
					Timestamp: migration.Timestamp,
					Name:      migration.Down,
					IsUp:      false,
				}
				m.On("Execute", expectedExec).Return(nil)
			}

			err := run.executeDownMigrations(v.down, v.downIDs)

			fs.AssertExpectations(t)
			m.AssertExpectations(t)

			if v.shouldReturnError {
				r.NotNil(err, "should return error")
			} else {
				r.Nil(err, "should not return error")
			}
		})
	}
}

func TestGetInDBDownMigrations(t *testing.T) {
	table := []struct {
		name   string
		border time.Time
		inDB   []int64
		result []int64
	}{
		{
			name:   "no migrations in database",
			border: time.Unix(20, 50),
			inDB:   []int64{},
			result: []int64{},
		},
		{
			name:   "downgrade all migrations",
			border: time.Unix(0, 0),
			inDB:   []int64{1, 2, 3, 4},
			result: []int64{4, 3, 2, 1},
		},
		{
			name:   "downgrade last 2 migrations",
			border: time.Unix(3, 0),
			inDB:   []int64{1, 2, 3, 1000, 1001},
			result: []int64{1001, 1000},
		},
	}

	for _, v := range table {
		t.Run(v.name, func(t *testing.T) {
			r := require.New(t)

			run := Run{}
			res := run.getInDBDownMigrations(v.inDB, v.border)

			r.Equal(res, v.result, "return values are not equal")
		})
	}
}
