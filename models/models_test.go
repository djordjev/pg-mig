package models

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var demoError = errors.New("demo error")

func TestCreateMetaTable(t *testing.T) {
	r := require.New(t)
	table := []struct {
		name  string
		query string
		err   error
	}{
		{
			name:  "successfully executes",
			query: fmt.Sprintf(createMetaTableQuery, tableName),
			err:   nil,
		},
		{
			name:  "returns error",
			query: fmt.Sprintf(createMetaTableQuery, tableName),
			err:   errors.New("some error"),
		},
	}

	for _, val := range table {
		t.Run(val.name, func(t *testing.T) {
			mockConnection := &mockedDBConnection{}
			mockConnection.On("Exec", mock.Anything, val.query, mock.Anything).
				Return(pgconn.CommandTag{}, val.err)

			m := ImplModels{Db: mockConnection}
			err := m.CreateMetaTable()

			mockConnection.AssertExpectations(t)

			if val.err == nil {
				r.NoError(err, "should not return error")
			} else {
				r.Error(err, "should return error")
			}
		})
	}

}

func TestGetMigrationsList(t *testing.T) {
	r := require.New(t)
	t1, _ := time.Parse(time.RFC3339, "2020-09-20T15:04:05Z")
	t2, _ := time.Parse(time.RFC3339, "2020-09-20T15:05:05Z")

	scanError := errors.New("scan error")

	table := []struct {
		queryError    error
		scanError     error
		queryRes      []time.Time
		expected      []int64
		name          string
		expectedError error
	}{
		{
			queryError:    nil,
			queryRes:      []time.Time{t1, t2},
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
			queryRes:      []time.Time{},
			expected:      []int64{},
			name:          "returns no results",
			expectedError: nil,
		},
		{
			queryError:    nil,
			queryRes:      []time.Time{t1, t2},
			expected:      []int64{},
			name:          "scan returns error",
			scanError:     scanError,
			expectedError: scanError,
		},
	}

	for _, val := range table {
		t.Run(val.name, func(ts *testing.T) {
			db := &mockedDBConnection{}
			rows := &rowsImpl{}

			db.On("Query", mock.Anything, fmt.Sprintf(getMigrationsListQuery, tableName), mock.Anything).
				Return(rows, val.queryError)

			rows.On("Close")
			rows.On("Scan", mock.Anything).Return(val.scanError)

			for range val.queryRes {
				rows.On("Next").Return(true).Once()
			}
			rows.On("Next").Return(false).Once()

			rows.scans = make([]interface{}, len(val.queryRes))
			for i := range val.queryRes {
				rows.scans[i] = val.queryRes[i]
			}

			m := ImplModels{Db: db}

			res, err := m.GetMigrationsList()

			if val.expectedError == nil {
				r.NoError(err)
			} else {
				r.Error(err)
			}

			r.Equal(res, val.expected)

		})

	}

}

func TestExecute(t *testing.T) {
	r := require.New(t)

	table := []struct {
		name             string
		executionContext ExecutionContext
		expectedMeta     string
		metaErr          error
		sqlErr           error
		commitErr        error
		returnError      error
	}{
		{
			name:             "executes up migration",
			executionContext: ExecutionContext{Timestamp: 123, Name: "demo_name", Sql: "sql", IsUp: true},
			expectedMeta:     fmt.Sprintf("insert into %s (ts) values ($1);", tableName),
			metaErr:          nil,
			sqlErr:           nil,
			commitErr:        nil,
			returnError:      nil,
		},
		{
			name:             "executes up migration meta table error",
			executionContext: ExecutionContext{Timestamp: 123, Name: "demo_name", Sql: "sql", IsUp: true},
			expectedMeta:     fmt.Sprintf("insert into %s (ts) values ($1);", tableName),
			metaErr:          errors.New("meta error"),
			sqlErr:           nil,
			commitErr:        nil,
			returnError:      errors.New("meta error"),
		},
		{
			name:             "executes up migration execution error",
			executionContext: ExecutionContext{Timestamp: 123, Name: "demo_name", Sql: "sql", IsUp: true},
			expectedMeta:     fmt.Sprintf("insert into %s (ts) values ($1);", tableName),
			metaErr:          nil,
			sqlErr:           errors.New("exec error"),
			commitErr:        nil,
			returnError:      errors.New("exec error"),
		},
		{
			name:             "executes up migration commit error",
			executionContext: ExecutionContext{Timestamp: 123, Name: "demo_name", Sql: "sql", IsUp: true},
			expectedMeta:     fmt.Sprintf("insert into %s (ts) values ($1);", tableName),
			metaErr:          nil,
			sqlErr:           nil,
			commitErr:        errors.New("commit error"),
			returnError:      nil,
		},
		{
			name:             "executes down migration",
			executionContext: ExecutionContext{Timestamp: 444, Name: "demo_name", Sql: "sql dn", IsUp: false},
			expectedMeta:     fmt.Sprintf("delete from %s where ts = $1", tableName),
			metaErr:          nil,
			sqlErr:           nil,
			commitErr:        nil,
			returnError:      nil,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			mockConn := mockedDBConnection{}
			tx := txImpl{}

			mockConn.On("Begin", mock.Anything).Return(&tx, nil)

			m := ImplModels{Db: &mockConn}

			// Calls exec on tx to update meta table
			ts := time.Unix(test.executionContext.Timestamp, 0)
			tx.On("Exec", mock.Anything, test.expectedMeta, []interface{}{ts}).
				Return(pgconn.CommandTag{}, test.metaErr).Once()

			// Calls exec with provided sql
			tx.On("Exec", mock.Anything, test.executionContext.Sql, mock.Anything).
				Return(pgconn.CommandTag{}, test.sqlErr).Once()

			// Calls commit
			tx.On("Commit", mock.Anything).Return(test.commitErr).Once()

			// Calls rollback
			tx.On("Rollback", mock.Anything).Return(nil)

			var err error
			if test.commitErr == nil {
				err = m.Execute(test.executionContext)
			} else {
				r.PanicsWithError(test.commitErr.Error(), func() {
					err = m.Execute(test.executionContext)
				})
			}

			if test.returnError != nil {
				r.Error(test.returnError)
			} else {
				r.NoError(err)
			}
		})
	}
}

func TestSquashMigrations(t *testing.T) {
	t.Skip("To implement")
}
