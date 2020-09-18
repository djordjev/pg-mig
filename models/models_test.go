package models

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
)

// Mocks
type MockedDBConnection struct {
	pgx.Conn
	err error
}

func (conn MockedDBConnection) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return nil, conn.err
}

func (conn MockedDBConnection) Close(_ context.Context) error {
	return nil
}

// Structs
var testModels struct {
	mockDBSuccess DBConnection
	mockDBError   DBConnection
}

func TestMain(m *testing.M) {
	testModels.mockDBSuccess = MockedDBConnection{err: nil}
	testModels.mockDBError = MockedDBConnection{err: errors.New("some err")}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestCreateMetaTable(t *testing.T) {
	m := Models{Db: testModels.mockDBSuccess}
	err := m.CreateMetaTable()
	if err != nil {
		t.Logf("Expected error to be nil but got %v", err)
		t.Fail()
	}
}

func TestCreateMetaTableFail(t *testing.T) {
	m := Models{Db: testModels.mockDBError}
	err := m.CreateMetaTable()
	if err == nil {
		t.Log("Expected to get error but got nil")
		t.Fail()
	}
}
