package repository

import (
	"context"
	"database/sql"
)

// DBTX is an interface that sql.DB and sql.Tx implement.
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

const (
	userTable     = "users"
	locationTable = "locations"
)
