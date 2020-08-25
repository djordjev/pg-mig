package models_test

import (
	"errors"
	"github.com/djordjev/pg-mig/testutils"
	"testing"

	"github.com/djordjev/pg-mig/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ModelsSuite struct {
	suite.Suite
	mockDB testutils.MockDBConnection
}

func (suite *ModelsSuite) SetupTest() {
	suite.mockDB = testutils.NewMockedDBConnection()
}

func (suite *ModelsSuite) TestCreateMetaTable() {
	r := suite.Require()

	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	m := models.Models{Db: &suite.mockDB}
	err := m.CreateMetaTable()
	r.Nil(err)
}

func (suite *ModelsSuite) TestCreateMetaTableFail() {
	r := suite.Require()

	suite.mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("err"))

	m := models.Models{Db: &suite.mockDB}
	err := m.CreateMetaTable()
	r.NotNil(err)
}

func TestModelsSuite(t *testing.T) {
	suite.Run(t, new(ModelsSuite))
}
