package testutils

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockDBConnection struct {
	mock.Mock
}

func (conn *MockDBConnection) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := conn.Called(ctx, sql, arguments)
	return nil, args.Error(1)
}

func NewMockedDBConnection() MockDBConnection {
	return MockDBConnection{}
}
