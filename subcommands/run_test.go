package subcommands

import (
	"github.com/djordjev/pg-mig/models"
	"testing"
	"time"
)

type mockedModels struct {
	models.Models
	migrationList      []int64
	migrationListError error
}

func (models mockedModels) GetMigrationsList() ([]int64, error) {
	return models.migrationList, models.migrationListError
}

func getNow() time.Time {
	now, _ := time.Parse(time.RFC3339, "2020-09-20T15:00:00Z")
	return now
}

func TestGetMigrationFiles(t *testing.T) {
	run := Run{
		CommandBase: CommandBase{
			Models: mockedModels{},
		},
		GetNow: getNow,
	}
}
