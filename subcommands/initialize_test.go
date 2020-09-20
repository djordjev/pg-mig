package subcommands

import (
	"context"
	"errors"
	"github.com/djordjev/pg-mig/filesystem"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"testing"

	"github.com/djordjev/pg-mig/models"
	"github.com/spf13/afero"
)

// Mocks
type MockedDBConnection struct {
	pgx.Conn
	err     error
	counter int8
}

func (conn MockedDBConnection) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	conn.counter++
	return nil, conn.err
}

func buildCommandBase(connection *MockedDBConnection) *CommandBase {
	cb := CommandBase{
		Models:     models.Models{Db: connection},
		Filesystem: afero.NewMemMapFs(),
		Config:     filesystem.Config{},
		Flags:      []string{},
	}

	return &cb
}

func TestCreateTableSuccess(t *testing.T) {
	commandBase := buildCommandBase(&MockedDBConnection{err: nil})
	initialize := Initialize{CommandBase: *commandBase}
	err := initialize.Run()

	if err != nil {
		t.Logf("Expected error to be nil but got: %v", err)
		t.Fail()
	}

}

func TestCreateTableError(t *testing.T) {
	commandBase := buildCommandBase(&MockedDBConnection{err: errors.New("some error")})

	initialize := Initialize{CommandBase: *commandBase}

	err := initialize.Run()

	if err == nil {
		t.Log("Expected to return error but got nil")
		t.Fail()
	}
}
