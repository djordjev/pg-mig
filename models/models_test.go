package models

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"os"
	"testing"
)

var testModels struct {
	mockDBSuccess DBConnection
	mockDBError DBConnection
}

type MockedDBConnectionSuccess struct {
	pgx.Conn
}

func (conn MockedDBConnectionSuccess) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

type MockedDBConnectionError struct {
	pgx.Conn
}

func (conn MockedDBConnectionError) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return nil, errors.New("some error")
}

func TestMain(m *testing.M) {
	testModels.mockDBSuccess = MockedDBConnectionSuccess{}
	testModels.mockDBError = MockedDBConnectionError{}

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
