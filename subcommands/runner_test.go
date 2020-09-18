package subcommands

import (
	"context"
	"fmt"
	"github.com/djordjev/pg-mig/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetSubcommand(t *testing.T) {
	table := []struct {
		runner   Runner
		hasError bool
		hasType  reflect.Type
	}{
		{runner: Runner{Subcommand: cmdInit}, hasError: false, hasType: reflect.TypeOf(&Initialize{})},
		{runner: Runner{Subcommand: "unknown"}, hasError: true, hasType: reflect.TypeOf(nil)},
	}

	for i := 0; i < len(table); i++ {
		test := table[i]

		cmd, err := test.runner.getSubcommand(&CommandBase{})

		returnedError := err != nil
		if test.hasError != returnedError {
			t.Logf("getSubcommand returns invalid value for error. Subcommand: %s, expected to have err %t, got %v", test.runner.Subcommand, test.hasError, returnedError)
			t.Fail()
		}

		returnedType := reflect.TypeOf(cmd)

		if returnedType != test.hasType {
			t.Logf("getSubcommand returned type %s but expected %s", returnedType.Name(), test.hasType.Name())
			t.Fail()
		}
	}
}

func TestCreateInitFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Log("Unable to get current working directory. Skipping test")
		t.Skip()
		return
	}

	table := []struct {
		flags             []string
		fs                afero.Fs
		configFileContent string
		hasError          bool
	}{
		{
			flags: []string{
				"-name=main_db",
				"-credentials=postgres:pg_pass",
			},
			fs:                afero.NewMemMapFs(),
			hasError:          false,
			configFileContent: fmt.Sprintf("{\"db_name\":\"main_db\",\"path\":\"%s\",\"db_url\":\"localhost\",\"credentials\":\"postgres:pg_pass\",\"port\":5432,\"ssl_mode\":\"disable\"}", wd),
		},
		{
			flags: []string{
				"-name=main_db2",
				"-credentials=pgs",
				"-path=some_path",
			},
			fs:                afero.NewMemMapFs(),
			hasError:          false,
			configFileContent: "{\"db_name\":\"main_db2\",\"path\":\"some_path\",\"db_url\":\"localhost\",\"credentials\":\"pgs\",\"port\":5432,\"ssl_mode\":\"disable\"}",
		},
		{
			flags: []string{
				"-name=main_db3",
				"-credentials=pgs11",
				"-path=some_path2",
				"-db=db",
				"-ssl=on",
				"-port=1111",
			},
			fs:                afero.NewMemMapFs(),
			hasError:          false,
			configFileContent: "{\"db_name\":\"main_db3\",\"path\":\"some_path2\",\"db_url\":\"db\",\"credentials\":\"pgs11\",\"port\":1111,\"ssl_mode\":\"on\"}",
		},
		{
			flags: []string{
				"-credentials=pgs",
			},
			fs:                afero.NewMemMapFs(),
			hasError:          true,
			configFileContent: "",
		},
		{
			flags: []string{
				"-name=dbname",
			},
			fs:                afero.NewMemMapFs(),
			hasError:          true,
			configFileContent: "",
		},
	}

	for i := 0; i < len(table); i++ {
		test := table[i]
		runner := Runner{Flags: test.flags, Filesystem: test.fs}

		err := runner.createInitFile()

		returnedError := err != nil
		if returnedError != test.hasError {
			t.Logf("createInitFile expected to return error %t but got %t", test.hasError, returnedError)
			t.Fail()
		}

		if returnedError {
			continue
		}

		bytesContent, err := afero.ReadFile(test.fs, "pgmig.config.json")
		if err != nil {
			t.Log("Can't read written config file")
			t.Fail()
			return
		}

		strContent := string(bytesContent)

		if test.configFileContent != strContent {
			t.Fail()
		}
	}
}

type connection struct {
	pgx.Conn
}

func (c connection) Close(_ context.Context) error {
	return nil
}

func (c connection) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func TestRun(t *testing.T) {
	connector := func(ctx context.Context, str string) (models.DBConnection, error) {
		return &connection{}, nil
	}

	runner := Runner{
		Filesystem: afero.NewMemMapFs(),
		Subcommand: cmdInit,
		Flags:      []string{"-name=main_db", "-credentials=postgres:pg_pass"},
		Connector:  connector,
	}

	err := runner.Run()

	if err != nil {
		t.Logf("Runner run function returned error %v", err)
		t.Fail()
	}

	bytesContent, err := afero.ReadFile(runner.Filesystem, "pgmig.config.json")
	content := string(bytesContent)

	hasDBName := strings.Index(content, "\"db_name\":\"main_db\"") > 0
	hasCredentials := strings.Index(content, "\"credentials\":\"postgres:pg_pass\"") > 0

	if !hasDBName {
		t.Log("Missing database name in config file")
		t.Fail()
	}

	if hasCredentials {
		t.Log("Missing credentials in configuration")
		t.Fail()
	}
}
