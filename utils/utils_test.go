package utils

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type UtilsSuite struct {
	suite.Suite
	configContent string
	Filesystem    afero.Fs
}

func (suite *UtilsSuite) SetupSuite() {
	suite.configContent = `{"db_name":"main_db","path":".","db_url":"localhost","credentials":"postgres:pg_pass","port":5432,"ssl_mode":"disable"}`
}

func (suite *UtilsSuite) SetupTest() {
	suite.Filesystem = afero.NewMemMapFs()
}

func (suite *UtilsSuite) TestLoadNoFile() {
	config := Config{Filesystem: suite.Filesystem}
	err := config.Load()

	suite.Require().NotNil(err)
}

func (suite *UtilsSuite) TestLoadWithFile() {
	config := Config{Filesystem: suite.Filesystem}

	afero.WriteFile(suite.Filesystem, configFileName, []byte(suite.configContent), 0644)

	err := config.Load()

	suite.Require().Nil(err)

	suite.Require().Equal(config.DbName, "main_db")
	suite.Require().Equal(config.Path, ".")
	suite.Require().Equal(config.DbURL, "localhost")
	suite.Require().Equal(config.Credentials, "postgres:pg_pass")
	suite.Require().Equal(config.Port, 5432)
	suite.Require().Equal(config.SSL, "disable")
}

func (suite *UtilsSuite) TestLoadInvalidFormat() {
	config := Config{Filesystem: suite.Filesystem}

	afero.WriteFile(suite.Filesystem, configFileName, []byte("not a json"), 0644)

	err := config.Load()

	suite.Require().NotNil(err)
}

func (suite *UtilsSuite) TestStoreSavesConfig() {
	config := Config{
		Credentials: "Credentials",
		DbName:      "DbName",
		DbURL:       "DbURL",
		Path:        "Path",
		Port:        1234,
		SSL:         "SSL",

		Filesystem: suite.Filesystem,
	}

	afero.WriteFile(suite.Filesystem, configFileName, []byte("wrong content"), 0777)

	err := config.Store()

	suite.Require().Nil(err)

	bytesContent, err := afero.ReadFile(suite.Filesystem, configFileName)
	if err != nil {
		suite.Fail("Can't read written config file")
		return
	}

	strContent := string(bytesContent)

	hasCredientials := strings.Index(strContent, "Credentials")
	hasDBName := strings.Index(strContent, "DbName")
	hasURL := strings.Index(strContent, "DbURL")
	hasPath := strings.Index(strContent, "Path")
	hasPort := strings.Index(strContent, "1234")
	hasSSL := strings.Index(strContent, "SSL")
	hasOldContent := strings.Index(strContent, "wrong content")

	suite.Require().False(hasOldContent >= 0)
	suite.Require().True(hasSSL >= 0)
	suite.Require().True(hasPort >= 0)
	suite.Require().True(hasPath >= 0)
	suite.Require().True(hasURL >= 0)
	suite.Require().True(hasDBName >= 0)
	suite.Require().True(hasCredientials >= 0)
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}
