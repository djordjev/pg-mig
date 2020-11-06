package subcommands

import "fmt"

// Help struct for help command
type Help struct{}

func (help *Help) Run() error {
	fmt.Println("Commands:")
	fmt.Println("init -> initializes pg-mig with a database to run migrations against")
	fmt.Println("add -> adds new migration files with current timestamp associated")
	fmt.Println("log -> prints available migrations in database and on filesystem")
	fmt.Println("run -> executes migrations for given time")
	fmt.Println("squash -> merges (squashes) multiple migrations into one")
	fmt.Println()
	fmt.Println("Note: for more info and flags run pg-mig command -help (for example pg-mig init -help)")
	return nil
}
