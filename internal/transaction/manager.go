package transaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(tx Transaction) error) error
}

type SQLXTransactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &SQLXTransactionManager{db: db}
}

func (m *SQLXTransactionManager) RunInTransaction(ctx context.Context, fn func(tx Transaction) error) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	sqlxTx := NewSQLXTransaction(tx)
	if err := fn(sqlxTx); err != nil {
		err = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		err = tx.Rollback()
		return err
	}

	return nil
}
