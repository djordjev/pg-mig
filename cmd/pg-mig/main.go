package main

import (
	"fmt"
	"github.com/djordjev/pg-mig/models"
	"os"

	"github.com/djordjev/pg-mig/subcommands"
	"github.com/spf13/afero"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Please run pg-mig help to see more info")
		return
	}

	runner := subcommands.Runner{
		Subcommand: os.Args[1],
		Flags:      os.Args[2:],
		Filesystem: afero.NewOsFs(),
		Connector:  models.BuildConnector,
	}

	err := runner.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}
