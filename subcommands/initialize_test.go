package subcommands_test

import (
	"errors"
	"github.com/djordjev/pg-mig/testutils"
	"testing"

	"github.com/djordjev/pg-mig/models"
	"github.com/djordjev/pg-mig/subcommands"
	"github.com/djordjev/pg-mig/utils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InitializeSuite struct {
	suite.Suite
	base   subcommands.CommandBase
	mockDB testutils.MockDBConnection
}

func (suite *InitializeSuite) SetupTest() {
	suite.mockDB = testutils.NewMockedDBConnection()

	suite.base = subcommands.CommandBase{
		Models:     models.Models{Db: &suite.mockDB},
		Filesystem: afero.NewMemMapFs(),
		Config:     utils.Config{},
		Flags:      []string{},
	}
}

func (suite *InitializeSuite) TestCreateTableSuccess() {
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	initialize := subcommands.Initialize{CommandBase: suite.base}
	err := initialize.Run()

	suite.Require().Nil(err)

	suite.mockDB.MethodCalled("Exec")
}

func (suite *InitializeSuite) TestCreateTableError() {
	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	initialize := subcommands.Initialize{CommandBase: suite.base}
	err := initialize.Run()

	suite.Require().NotNil(err)

	suite.mockDB.MethodCalled("Exec")
}

func TestInitializeSuite(t *testing.T) {
	suite.Run(t, new(InitializeSuite))
}
