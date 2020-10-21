package models

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"reflect"
)

type mockedDBConnection struct {
	mock.Mock
}

func (m *mockedDBConnection) Begin(ctx context.Context) (pgx.Tx, error) {
	c := m.Called(ctx)
	return c.Get(0).(*txImpl), c.Error(1)
}

func (m *mockedDBConnection) Close(ctx context.Context) error {
	c := m.Called(ctx)
	return c.Error(0)
}

func (m *mockedDBConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	c := m.Called(ctx, sql, arguments)
	return c.Get(0).(pgconn.CommandTag), c.Error(1)
}

func (m *mockedDBConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	c := m.Called(ctx, sql, args)
	return c.Get(0).(*rowsImpl), c.Error(1)
}

type rowsImpl struct {
	mock.Mock
	scans []interface{}
	cnt   int
}

func (r *rowsImpl) Close() {
	r.Called()
}

func (r *rowsImpl) Err() error {
	c := r.Called()
	return c.Error(0)
}

func (r *rowsImpl) CommandTag() pgconn.CommandTag {
	c := r.Called()
	return c.Get(0).(pgconn.CommandTag)
}

func (r *rowsImpl) FieldDescriptions() []pgproto3.FieldDescription {
	c := r.Called()
	return c.Get(0).([]pgproto3.FieldDescription)
}

func (r *rowsImpl) Next() bool {
	c := r.Called()
	return c.Bool(0)
}

func (r *rowsImpl) Scan(dest ...interface{}) error {
	c := r.Called(dest)
	cs := r.scans[r.cnt]
	defer func() { r.cnt++ }()

	vres := reflect.ValueOf(cs)
	reflect.ValueOf(dest[0]).Elem().Set(vres)

	return c.Error(0)
}

func (r *rowsImpl) Values() ([]interface{}, error) {
	c := r.Called()
	return c.Get(0).([]interface{}), c.Error(1)
}

func (r *rowsImpl) RawValues() [][]byte {
	c := r.Called()
	return c.Get(0).([][]byte)
}

type txImpl struct {
	mock.Mock
}

func (t *txImpl) Begin(ctx context.Context) (pgx.Tx, error) {
	c := t.Called(ctx)
	return c.Get(0).(pgx.Tx), c.Error(1)
}

func (t *txImpl) Commit(ctx context.Context) error {
	c := t.Called(ctx)
	return c.Error(0)
}

func (t *txImpl) Rollback(ctx context.Context) error {
	c := t.Called(ctx)
	return c.Error(0)
}

func (t *txImpl) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	c := t.Called(ctx, tableName, columnNames, rowSrc)
	return c.Get(0).(int64), c.Error(1)
}

func (t *txImpl) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	c := t.Called(ctx, b)
	return c.Get(0).(pgx.BatchResults)
}

func (t *txImpl) LargeObjects() pgx.LargeObjects {
	c := t.Called()
	return c.Get(0).(pgx.LargeObjects)
}

func (t *txImpl) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	c := t.Called(ctx, name, sql)
	return c.Get(0).(*pgconn.StatementDescription), c.Error(1)
}

func (t *txImpl) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	c := t.Called(ctx, sql, arguments)
	return c.Get(0).(pgconn.CommandTag), c.Error(1)
}

func (t *txImpl) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	c := t.Called(ctx, sql, args)
	return c.Get(0).(*rowsImpl), c.Error(1)
}

func (t *txImpl) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	c := t.Called(ctx, sql, args)
	return c.Get(0).(*rowsImpl)
}

func (t *txImpl) Conn() *pgx.Conn {
	c := t.Called()
	return c.Get(0).(*pgx.Conn)
}
