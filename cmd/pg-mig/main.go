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
	printer := subcommands.ImplPrinter{NoColor: true}

	if len(os.Args) < 2 {
		printer.PrintError("Missing command. Please run pg-mig help to see more info")
		os.Exit(1)
	}

	runner := subcommands.Runner{
		Subcommand: os.Args[1],
		Flags:      os.Args[2:],
		Fs:         &filesystem.ImplFilesystem{Fs: afero.NewOsFs(), GetNow: time.Now},
		Connector:  models.BuildConnector,
		Timer:      timer.Timer{Now: time.Now},
		Printer:    &printer,
	}

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("execution error: ", e)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}()

	err := runner.Run()
	if err != nil {
		printer.PrintError(err.Error())
		os.Exit(1)
	}
}
