package subcommands

import (
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/spf13/afero"
	"testing"
	"time"
)

func TestAddRun(t *testing.T) {
	buildGetTime := func(t string) filesystem.TimeGetter {
		return func() time.Time {
			time, _ := time.Parse(time.RFC3339, t)
			return time
		}
	}

	fs := afero.NewMemMapFs()

	table := []struct {
		fs         afero.Fs
		timeGetter filesystem.TimeGetter
		flags      []string
		name       string
	}{
		{
			fs:         fs,
			timeGetter: buildGetTime("2020-09-20T15:04:05Z"),
			flags:      []string{},
			name:       "mig_1600614245",
		},
		{
			fs:         fs,
			timeGetter: buildGetTime("2020-09-20T15:05:05+07:00"),
			flags:      []string{},
			name:       "mig_1600589105",
		},
		{
			fs:         fs,
			timeGetter: buildGetTime("2020-09-20T15:09:05-02:00"),
			flags:      []string{"-name=random"},
			name:       "mig_1600621745_random",
		},
	}

	for i := 0; i < len(table); i++ {
		current := table[i]
		add := Add{
			CommandBase{
				Filesystem: filesystem.Filesystem{Fs: fs, GetNow: current.timeGetter},
				Flags:      current.flags,
			},
		}

		add.Run()

		existsUp, _ := afero.Exists(fs, fmt.Sprintf("%s_up.sql", current.name))
		if !existsUp {
			t.Logf("Up file for %s is missing", current.name)
			t.Fail()
		}

		existsDown, _ := afero.Exists(fs, fmt.Sprintf("%s_down.sql", current.name))
		if !existsDown {
			t.Logf("Down file for %s is missing", current.name)
			t.Fail()
		}
	}
}
