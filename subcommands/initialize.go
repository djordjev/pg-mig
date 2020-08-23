package subcommands

// Initialize structure for init command
type Initialize struct {
	CommandBase
}

// Run function creates meta table in db
func (init *Initialize) Run() error {
	err := init.Models.CreateMetaTable()
	if err != nil {
		return err
	}

	return nil
}
