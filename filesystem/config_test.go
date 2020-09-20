package filesystem

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

var validContent = `{"db_name":"main_db","path":".","db_url":"localhost","credentials":"postgres:pg_pass","port":5432,"ssl_mode":"disable"}`

func TestLoadNoFile(t *testing.T) {
	config := Config{Filesystem: afero.NewMemMapFs()}
	err := config.Load()

	if err == nil {
		t.Logf("TestLoadNoFile fails because error was nil")
		t.Fail()
	}
}

func TestLoadWithFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	config := Config{Filesystem: fs}

	err := afero.WriteFile(fs, configFileName, []byte(validContent), 0644)
	if err != nil {
		t.Log("Can't write test file in TestLoadWithFile")
		t.Fail()
		return
	}

	err = config.Load()

	if err != nil {
		t.Logf("Expected err to be nil but got %v", err)
		t.Fail()
	}

	if config.DbName != "main_db" {
		t.Logf("Expected DbName to be 'main_db' but got %s", config.DbName)
		t.Fail()
	}

	if config.Path != "." {
		t.Logf("Expected Path to be '.' but got %s", config.Path)
		t.Fail()
	}

	if config.DbURL != "localhost" {
		t.Logf("Expected DbURL to be localhost but got %s", config.DbURL)
		t.Fail()
	}

	if config.Credentials != "postgres:pg_pass" {
		t.Logf("Expected credentials to be postgres:pg_pass but got %s", config.Credentials)
		t.Fail()
	}

	if config.Port != 5432 {
		t.Logf("Expected port to be 5432 but got %d", config.Port)
		t.Fail()
	}

	if config.SSL != "disable" {
		t.Logf("Expected ssl to be disable but got %s", config.SSL)
		t.Fail()
	}
}

func TestLoadInvalidFormat(t *testing.T) {
	fs := afero.NewMemMapFs()

	config := Config{Filesystem: fs}

	err := afero.WriteFile(fs, configFileName, []byte("not a json"), 0644)
	if err != nil {
		t.Log("Can't write test file in TestLoadInvalidFormat")
		t.Fail()
		return
	}

	err = config.Load()

	if err == nil {
		t.Log("Expected to return error when format is invalid")
		t.Fail()
	}
}

func TestStoreSavesConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	config := Config{
		Credentials: "Credentials",
		DbName:      "DbName",
		DbURL:       "DbURL",
		Path:        "Path",
		Port:        1234,
		SSL:         "SSL",

		Filesystem: fs,
	}

	err := afero.WriteFile(fs, configFileName, []byte("wrong content"), 0777)
	if err != nil {
		t.Log("Can't write test file in TestStoreSavesConfig")
		t.Fail()
		return
	}

	err = config.Store()

	if err != nil {
		t.Logf("Expected to get nil for error but got %v", err)
		t.Fail()
	}

	bytesContent, err := afero.ReadFile(fs, configFileName)
	if err != nil {
		t.Log("Can't read written config file")
		t.Fail()
		return
	}

	strContent := string(bytesContent)

	hasCredentials := strings.Index(strContent, "Credentials")
	hasDBName := strings.Index(strContent, "DbName")
	hasURL := strings.Index(strContent, "DbURL")
	hasPath := strings.Index(strContent, "Path")
	hasPort := strings.Index(strContent, "1234")
	hasSSL := strings.Index(strContent, "SSL")
	hasOldContent := strings.Index(strContent, "wrong content")

	if hasOldContent >= 0 {
		t.Log("Found old content after rewriting config")
		t.Fail()
	}

	if hasSSL < 0 {
		t.Log("Not fount SSL")
		t.Fail()
	}

	if hasPort < 0 {
		t.Log("Not found port")
		t.Fail()
	}

	if hasPath < 0 {
		t.Log("Not found path")
		t.Fail()
	}

	if hasURL < 0 {
		t.Log("Not found URL")
		t.Fail()
	}

	if hasDBName < 0 {
		t.Log("Not found DB Name")
		t.Fail()
	}

	if hasCredentials < 0 {
		t.Log("Not found credentials")
		t.Fail()
	}
}

func TestGetConnectionString(t *testing.T) {
	fs := afero.NewMemMapFs()

	config := Config{
		Credentials: "Credentials",
		DbName:      "DbName",
		DbURL:       "DbURL",
		Path:        "Path",
		Port:        1234,
		SSL:         "SSL",

		Filesystem: fs,
	}

	connectionString, err := config.GetConnectionString()

	if err != nil {
		t.Logf("Expected err to be nil but got %v", err)
		t.Fail()
	}

	if connectionString != "postgres://Credentials@DbURL:1234/DbName?sslmode=SSL" {
		t.Logf("Got invalid connection string %s", connectionString)
		t.Fail()
	}

	config.Credentials = ""
	connectionString, err = config.GetConnectionString()

	if connectionString != "" {
		t.Logf("Expected to get empty string but got %s", connectionString)
		t.Fail()
	}

	if err == nil {
		t.Log("Expected to get error but err is nil")
		t.Fail()
	}
}
