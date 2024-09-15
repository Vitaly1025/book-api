package transaction

import (
	"context"
	"database/sql"
)

type Transaction interface {
	Commit() error
	Rollback() error

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
