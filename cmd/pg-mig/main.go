package main

import (
	"fmt"
	"os"

	"github.com/djordjev/pg-mig/commands"
)

const cmdInit = "init"
const cmdHelp = "help"

type command = func() error

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Please run pg-mig help to see more info")
		return
	}

	cmd := getCommand()
	if cmd != nil {
		err := cmd()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(fmt.Sprintf("Invalid argument %s", os.Args[1]))
	}
}

func getCommand() command {
	switch os.Args[1] {
	case cmdInit:
		{
			return commands.Initialize
		}

	case cmdHelp:
		{
			return commands.Help
		}
	}

	return nil
}
