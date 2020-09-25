package models

import (
	"errors"
	"testing"
	"time"
)

func TestCreateMetaTable(t *testing.T) {
	m := Models{Db: &mockedDBConnection{execError: nil}}
	err := m.CreateMetaTable()
	if err != nil {
		t.Logf("Expected error to be nil but got %v", err)
		t.Fail()
	}
}

func TestCreateMetaTableFail(t *testing.T) {
	m := Models{Db: &mockedDBConnection{execError: errors.New("some err")}}
	err := m.CreateMetaTable()
	if err == nil {
		t.Log("Expected to get error but got nil")
		t.Fail()
	}
}

func TestGetMigrationsList(t *testing.T) {
	var queryRes [][]interface{}
	queryRes = make([][]interface{}, 2, 2)

	t1, _ := time.Parse(time.RFC3339, "2020-09-20T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2020-09-20T15:05:05Z")

	queryRes[0] = []interface{}{t1}
	queryRes[1] = []interface{}{t2}

	db := &mockedDBConnection{queryRes: queryRes}
	m := Models{Db: db}

	res, err := m.GetMigrationsList()
	if err != nil {
		t.Logf("GetMigrationsList returned error %v", err)
		t.Fail()
	}

	if res[0] != t1.Unix() || res[1] != t2.Unix() {
		t.Fail()
	}
}
