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

func (m *mockedModels) SquashMigrations(i time.Time, i2 time.Time, i3 int64) error {
	c := m.Called(i, i2, i3)
	return c.Error(0)
}

func (m *mockedModels) CreateMetaTable() error {
	return m.createMetaTableError
}

func (m *mockedModels) GetMigrationsList() ([]int64, error) {
	c := m.Called()
	return c.Get(0).([]int64), c.Error(1)
}

func (m *mockedModels) Execute(executionContext models.ExecutionContext) error {
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
}

func (m *mockedFilesystem) Squash(list filesystem.MigrationFileList) error {
	c := m.Called(list)
	return c.Error(0)
}

func (m *mockedFilesystem) DeleteMigrationFiles(list filesystem.MigrationFileList) error {
	c := m.Called(list)
	return c.Error(0)
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

func (m *mockedFilesystem) GetFileTimestamps(t1 time.Time, t2 time.Time) (filesystem.MigrationFileList, error) {
	args := m.Called(t1, t2)
	return args.Get(0).(filesystem.MigrationFileList), args.Error(1)
}

type mockedPrinter struct {
	mock.Mock
}

func (m *mockedPrinter) PrintUpMigration(text string) {
	m.Called(text)
}

func (m *mockedPrinter) PrintDownMigration(text string) {
	m.Called(text)
}

func (m *mockedPrinter) PrintError(text string) {
	m.Called(text)
}

func (m *mockedPrinter) PrintSuccess(text string) {
	m.Called(text)
}

func (m *mockedPrinter) PrintMigrations(date string, onFS string, inDB string) {
	m.Called(date, onFS, inDB)
}

func (m *mockedPrinter) SetNoColor(color bool) {
	m.Called(color)
}
