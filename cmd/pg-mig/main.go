package main

import (
	"fmt"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/timer"
	"github.com/spf13/afero"
	"os"
	"time"

	"github.com/djordjev/pg-mig/subcommands"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Please run pg-mig help to see more info")
		return
	}

	printer := subcommands.ImplPrinter{NoColor: true}

	runner := subcommands.Runner{
		Subcommand: os.Args[1],
		Flags:      os.Args[2:],
		Fs:         &filesystem.ImplFilesystem{Fs: afero.NewOsFs(), GetNow: time.Now},
		Connector:  models.BuildConnector,
		Timer:      timer.Timer{Now: time.Now},
		Printer:    &printer,
	}

	err := runner.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}
