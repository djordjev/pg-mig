package filesystem

import (
	"time"

	"github.com/spf13/afero"
)

type ImplFilesystem struct {
	Fs             afero.Fs
	GetNow         func() time.Time
	ConfigDir      string
	ExternalConfig *Config
}

type TimeGetter func() time.Time

type Direction int

const (
	DirectionUp = iota
	DirectionDown
)

func (d Direction) String() string {
	names := [...]string{"DirectionUp", "DirectionDown"}
	if int(d) < len(names) {
		return names[d]
	}

	return "UnknownDirection"
}

type Filesystem interface {
	StoreConfig(config Config) error
	LoadConfig() (Config, error)
	CreateMigrationFile(string, string) error
	ReadMigrationContent(MigrationFile, Direction, Config) (string, error)
	GetFileTimestamps(time.Time, time.Time) (MigrationFileList, error)
	Squash(MigrationFileList) error
}
