package subcommands

import (
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"time"
)

// Run structure for run command
type Run struct {
	CommandBase
}

// Run executes up/down migrations
func (run *Run) Run() error {
	files, err := run.Filesystem.GetFiles(time.Time{}, time.Now(), filesystem.DirectionUp)
	if err != nil {
		return err
	}

	fmt.Println(files)

	return nil
}
