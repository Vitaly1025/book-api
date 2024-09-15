package transaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func NewSQLXTransaction(tx *sqlx.Tx) Transaction {
	return &SQLXTransaction{tx: tx}
}

type SQLXTransaction struct {
	tx *sqlx.Tx
}

func (t *SQLXTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLXTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *SQLXTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *SQLXTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}
