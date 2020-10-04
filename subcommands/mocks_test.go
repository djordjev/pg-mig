package subcommands

import (
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/djordjev/pg-mig/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type mockedModels struct {
	mock.Mock
	createMetaTableError   error
	getMigrationsListError error
	getMigrationsListRes   []int64
	executeError           error
}

func (m mockedModels) CreateMetaTable() error {
	return m.createMetaTableError
}

func (m mockedModels) GetMigrationsList() ([]int64, error) {
	return m.getMigrationsListRes, m.getMigrationsListError
}

func (m mockedModels) Execute(executionContext models.ExecutionContext) error {
	c := m.Called(executionContext)
	return c.Error(0)
}

type mockedFilesystem struct {
	mock.Mock
	storeConfigError          error
	loadConfigError           error
	loadConfigConfig          filesystem.Config
	createMigrationFileError  error
	readMigrationContentRes   string
	readMigrationContentError error
	getFileTimestampsRes      []filesystem.MigrationFileList
	getFileTimestampsError    error
	getFileTimestampsResIter  int
}

func (m *mockedFilesystem) StoreConfig(_ filesystem.Config) error {
	return m.storeConfigError
}

func (m *mockedFilesystem) LoadConfig() (filesystem.Config, error) {
	return m.loadConfigConfig, m.loadConfigError
}

func (m *mockedFilesystem) CreateMigrationFile(_ string, _ string) error {
	return m.createMigrationFileError
}

func (m *mockedFilesystem) ReadMigrationContent(file filesystem.MigrationFile, direction filesystem.Direction, config filesystem.Config) (string, error) {
	c := m.Called(file, direction, config)
	return c.String(0), c.Error(1)
}

func (m *mockedFilesystem) GetFileTimestamps(_ time.Time, _ time.Time) (filesystem.MigrationFileList, error) {
	if m.getFileTimestampsError != nil {
		return nil, m.getFileTimestampsError
	}

	val := m.getFileTimestampsRes[m.getFileTimestampsResIter]
	m.getFileTimestampsResIter++

	return val, m.getFileTimestampsError
}
