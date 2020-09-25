package models

import (
	"context"
	"github.com/jackc/pgconn"
	pgproto32 "github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4"
	"reflect"
)

type mockedDBConnection struct {
	pgx.Conn
	execError  error
	queryError error
	queryRes   [][]interface{}
}

func (conn mockedDBConnection) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return nil, conn.execError
}

func (conn mockedDBConnection) Query(_ context.Context, query string, _ ...interface{}) (pgx.Rows, error) {
	if conn.queryError != nil {
		return nil, conn.queryError
	}

	return &rowsImpl{conn: conn}, nil
}

type rowsImpl struct {
	conn mockedDBConnection
	cnt  int
}

// pgx.Rows interface implementation
func (r *rowsImpl) Next() bool {
	res := r.cnt < len(r.conn.queryRes)
	r.cnt++
	return res
}

func (r *rowsImpl) Scan(dest ...interface{}) error {
	pos := r.cnt - 1
	current := r.conn.queryRes[pos]

	for i := 0; i < len(dest); i++ {
		c := current[i]
		d := dest[i]

		vres := reflect.ValueOf(c)
		reflect.ValueOf(d).Elem().Set(vres)
	}

	return nil
}

func (r *rowsImpl) Close()                                          {}
func (r *rowsImpl) Err() error                                      { panic("implement me") }
func (r *rowsImpl) CommandTag() pgconn.CommandTag                   { panic("implement me") }
func (r *rowsImpl) FieldDescriptions() []pgproto32.FieldDescription { panic("implement me") }
func (r *rowsImpl) Values() ([]interface{}, error)                  { panic("implement me") }
func (r *rowsImpl) RawValues() [][]byte                             { panic("implement me") }
