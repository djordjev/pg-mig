package filesystem

import "github.com/spf13/afero"

type Filesystem struct {
	Fs afero.Fs
}
