package filesystem

import (
	"github.com/spf13/afero"
	"time"
)

type Filesystem struct {
	Fs     afero.Fs
	GetNow TimeGetter
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
