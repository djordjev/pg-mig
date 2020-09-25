package subcommands

import "fmt"

// Run structure for run command
type Run struct {
	CommandBase
}

// Run executes up/down migrations
func (run *Run) Run() error {

	// TODO check file formats and matching down files
	inDBMigs, err := run.Models.GetMigrationsList()
	if err != nil {
		return err
	}

	fmt.Println(inDBMigs)

	return nil
}
