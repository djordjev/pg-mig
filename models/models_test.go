package models

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

var demoError = errors.New("demo error")

func TestCreateMetaTable(t *testing.T) {
	table := []struct {
		name      string
		execError error
	}{
		{
			name:      "successfully executes",
			execError: nil,
		},
		{
			name:      "returns error",
			execError: errors.New("some error"),
		},
	}

	for _, val := range table {
		t.Run(val.name, func(t *testing.T) {
			m := Models{Db: &mockedDBConnection{execError: val.execError}}
			err := m.CreateMetaTable()
			if err != val.execError {
				t.Fail()
			}
		})
	}

}

func TestGetMigrationsList(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2020-09-20T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2020-09-20T15:05:05Z")

	table := []struct {
		queryError    error
		scanError     error
		queryRes      [][]interface{}
		expected      []int64
		name          string
		expectedError error
	}{
		{
			queryError:    nil,
			queryRes:      [][]interface{}{{t1}, {t2}},
			expected:      []int64{t1.Unix(), t2.Unix()},
			name:          "two migrations",
			expectedError: nil,
		},
		{
			queryError:    demoError,
			queryRes:      nil,
			expected:      nil,
			name:          "throws an error",
			expectedError: demoError,
		},
		{
			queryError:    nil,
			queryRes:      [][]interface{}{},
			expected:      []int64{},
			name:          "returns no results",
			expectedError: nil,
		},
		{
			queryError:    nil,
			queryRes:      [][]interface{}{{t1}, {t2}},
			expected:      []int64{},
			name:          "scan returns error",
			scanError:     scanError,
			expectedError: scanError,
		},
	}

	for _, val := range table {
		t.Run(val.name, func(ts *testing.T) {
			db := &mockedDBConnection{
				queryRes:   val.queryRes,
				queryError: val.queryError,
				scanErr:    val.scanError,
			}
			m := Models{Db: db}

			res, err := m.GetMigrationsList()

			if err != val.expectedError {
				t.Logf("got error when executing GetMigrationsList %v", err)
				t.Fail()
			}

			if !reflect.DeepEqual(res, val.expected) {
				t.Logf("return value is not the same as expected: %v, %v", val.expected, res)
				t.Fail()
			}
		})

	}

}
