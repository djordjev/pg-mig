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
