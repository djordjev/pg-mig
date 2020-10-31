package filesystem

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateFileWithLocation(t *testing.T) {
	fs := afero.NewMemMapFs()
	fsystem := &ImplFilesystem{Fs: fs}

	err := fsystem.CreateMigrationFile("test_name.sql", "./test/")
	if err != nil {
		t.Logf("Unable to create migration file %v", err)
		t.Fail()
	}

	exists, err := afero.Exists(fs, "test/test_name.sql")
	if err != nil {
		t.Logf("Error when creating a file %v", err)
		t.Fail()
	}

	if !exists {
		t.Log("Created file does not exist on FS")
		t.Fail()
	}

	exists, err = afero.Exists(fs, "test_name.sql")
	if err != nil {
		t.Logf("Error when checking created file %v", err)
		t.Fail()
	}

	if exists {
		t.Log("Created file ignored location")
		t.Fail()
	}
}

func TestCreateFileNoLocation(t *testing.T) {
	fs := afero.NewMemMapFs()
	fsystem := &ImplFilesystem{Fs: fs}

	err := fsystem.CreateMigrationFile("no_loc_name.sql", "")
	if err != nil {
		t.Logf("Unable to create migration file %v", err)
		t.Fail()
	}

	exists, err := afero.Exists(fs, "no_loc_name.sql")
	if err != nil {
		t.Logf("Error when creating a file %v", err)
		t.Fail()
	}

	if !exists {
		t.Log("Created file does not exist on FS")
		t.Fail()
	}
}

func TestGetFileTimestamps(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := []struct {
		time string
		name string
	}{
		{
			time: "2020-09-20T12:00:05Z",
		},
		{
			time: "2020-09-20T13:00:05Z",
			name: "_demoname",
		},
		{
			time: "2020-09-20T14:00:05Z",
			name: "_testname",
		},
		{
			time: "2020-09-20T15:00:05Z",
		},
	}

	afero.WriteFile(fs, configFileName, []byte(validContent), 0644)

	for index, file := range files {
		time, _ := time.Parse(time.RFC3339, file.time)

		upFN := fmt.Sprintf("mig_%d%s_up.sql", time.Unix(), file.name)
		dnFN := fmt.Sprintf("mig_%d%s_down.sql", time.Unix(), file.name)

		afero.WriteFile(fs, upFN, []byte("demo up content"), 0644)
		afero.WriteFile(fs, dnFN, []byte("demo down content"), 0644)

		afero.WriteFile(fs, fmt.Sprintf("random_file_%d", index), []byte("rand"), 0644)
	}

	table := []struct {
		name   string
		from   string
		to     string
		result []struct {
			up   string
			down string
		}
	}{
		{
			name: "returns first file",
			from: "2019-09-20T12:00:05Z",
			to:   "2020-09-20T12:01:05Z",
			result: []struct {
				up   string
				down string
			}{{up: "mig_1600603205_up.sql", down: "mig_1600603205_down.sql"}},
		},
		{
			name: "returns all files",
			from: "2019-09-20T12:00:05Z",
			to:   "2021-09-20T12:01:05Z",
			result: []struct {
				up   string
				down string
			}{
				{up: "mig_1600603205_up.sql", down: "mig_1600603205_down.sql"},
				{up: "mig_1600606805_demoname_up.sql", down: "mig_1600606805_demoname_down.sql"},
				{up: "mig_1600610405_testname_up.sql", down: "mig_1600610405_testname_down.sql"},
				{up: "mig_1600614005_up.sql", down: "mig_1600614005_down.sql"},
			},
		},
		{
			name: "returns first two files",
			from: "2020-01-20T12:01:05Z",
			to:   "2020-09-20T13:00:05Z",
			result: []struct {
				up   string
				down string
			}{
				{up: "mig_1600603205_up.sql", down: "mig_1600603205_down.sql"},
				{up: "mig_1600606805_demoname_up.sql", down: "mig_1600606805_demoname_down.sql"},
			},
		},
		{
			name: "returns last 2 files",
			from: "2020-09-20T13:00:05Z",
			to:   "2021-09-20T12:01:05Z",
			result: []struct {
				up   string
				down string
			}{
				{up: "mig_1600610405_testname_up.sql", down: "mig_1600610405_testname_down.sql"},
				{up: "mig_1600614005_up.sql", down: "mig_1600614005_down.sql"},
			},
		},
	}

	for _, val := range table {
		t.Run(val.name, func(t *testing.T) {
			fsystem := &ImplFilesystem{Fs: fs}
			t1, _ := time.Parse(time.RFC3339, val.from)
			t2, _ := time.Parse(time.RFC3339, val.to)

			res, err := fsystem.GetFileTimestamps(t1, t2)

			if err != nil {
				t.Fail()
			}

			for k, v := range res {
				if v.Up != val.result[k].up {
					t.Fail()
				}

				if v.Down != val.result[k].down {
					t.Fail()
				}
			}
		})
	}

}

func TestReadMigrationContent(t *testing.T) {
	fs := afero.NewMemMapFs()
	fsystem := &ImplFilesystem{Fs: fs}

	file := MigrationFile{Up: "mig_123_up.sql", Down: "mig_123_down.sql", Timestamp: 123}
	config := Config{}

	afero.WriteFile(fs, "mig_123_up.sql", []byte("up mig"), 0644)
	afero.WriteFile(fs, "mig_123_down.sql", []byte("down mig"), 0644)

	content, err := fsystem.ReadMigrationContent(file, DirectionUp, config)
	if err != nil {
		t.Fail()
	}
	if content != "up mig" {
		t.Fail()
	}

	content, err = fsystem.ReadMigrationContent(file, DirectionDown, config)
	if err != nil {
		t.Fail()
	}
	if content != "down mig" {
		t.Fail()
	}

	config.Path = "/dasdas"
	_, err = fsystem.ReadMigrationContent(file, DirectionDown, config)
	if err == nil {
		t.Fail()
	}
}

func TestSquash(t *testing.T) {
	r := require.New(t)
	table := []struct {
		name           string
		files          MigrationFileList
		returnError    bool
		deletedFiles   []string
		remainingFiles []string
		upFilename     string
		upContent      string
		downFilename   string
		downContent    string
	}{
		{
			name: "successfully squashes files all files",
			files: MigrationFileList{
				MigrationFile{Timestamp: 1, Up: "mig_1_up.sql", Down: "mig_1_down.sql"},
				MigrationFile{Timestamp: 2, Up: "mig_2_up.sql", Down: "mig_2_down.sql"},
				MigrationFile{Timestamp: 3, Up: "mig_3_up.sql", Down: "mig_3_down.sql"},
			},
			deletedFiles:   []string{"mig_1_up.sql", "mig_2_up.sql", "mig_3_up.sql", "mig_1_down.sql", "mig_2_down.sql", "mig_3_down.sql"},
			remainingFiles: []string{"mig_3_squashed_up.sql", "mig_3_squashed_down.sql"},
			upFilename:     "mig_3_squashed_up.sql",
			upContent:      "-- migration 1 UP\nup mig 1\n-- migration 2 UP\nup mig 2\n-- migration 3 UP\nup mig 3\n",
			downFilename:   "mig_3_squashed_down.sql",
			downContent:    "-- migration 3 DOWN\ndown mig 3\n-- migration 2 DOWN\ndown mig 2\n-- migration 1 DOWN\ndown mig 1\n",
		},
		{
			name: "successfully squashes first two files",
			files: MigrationFileList{
				MigrationFile{Timestamp: 1, Up: "mig_1_up.sql", Down: "mig_1_down.sql"},
				MigrationFile{Timestamp: 2, Up: "mig_2_up.sql", Down: "mig_2_down.sql"},
			},
			deletedFiles:   []string{"mig_1_up.sql", "mig_2_up.sql", "mig_1_down.sql", "mig_2_down.sql"},
			remainingFiles: []string{"mig_2_squashed_up.sql", "mig_2_squashed_down.sql", "mig_3_up.sql", "mig_3_down.sql"},
			upFilename:     "mig_2_squashed_up.sql",
			upContent:      "-- migration 1 UP\nup mig 1\n-- migration 2 UP\nup mig 2\n",
			downFilename:   "mig_2_squashed_down.sql",
			downContent:    "-- migration 2 DOWN\ndown mig 2\n-- migration 1 DOWN\ndown mig 1\n",
		},
		{
			name:           "returns error",
			files:          MigrationFileList{},
			deletedFiles:   []string{},
			remainingFiles: []string{"mig_1_up.sql", "mig_2_up.sql", "mig_3_up.sql", "mig_1_down.sql", "mig_2_down.sql", "mig_3_down.sql"},
			returnError:    true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			fsystem := &ImplFilesystem{Fs: fs}

			afero.WriteFile(fs, "mig_1_up.sql", []byte("up mig 1"), 0644)
			afero.WriteFile(fs, "mig_1_down.sql", []byte("down mig 1"), 0644)

			afero.WriteFile(fs, "mig_2_up.sql", []byte("up mig 2"), 0644)
			afero.WriteFile(fs, "mig_2_down.sql", []byte("down mig 2"), 0644)

			afero.WriteFile(fs, "mig_3_up.sql", []byte("up mig 3"), 0644)
			afero.WriteFile(fs, "mig_3_down.sql", []byte("down mig 3"), 0644)

			afero.WriteFile(fs, configFileName, []byte(validContent), 0644)

			err := fsystem.Squash(test.files)

			for _, f := range test.deletedFiles {
				exists, _ := afero.Exists(fs, f)
				r.False(exists)
			}

			for _, f := range test.remainingFiles {
				exists, _ := afero.Exists(fs, f)
				r.True(exists)
			}

			if test.upFilename != "" {
				content, _ := afero.ReadFile(fs, test.upFilename)
				r.Equal(string(content), test.upContent)
			}

			if test.downFilename != "" {
				content, _ := afero.ReadFile(fs, test.downFilename)
				r.Equal(string(content), test.downContent)
			}

			if test.returnError {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
