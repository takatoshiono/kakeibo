package database

import (
	"context"
	"database/sql"
	"fmt"
)

// DB is an interface for DB and Tx of database/sql
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Transaction manages *sql.DB and current *sql.Tx
type Transaction struct {
	db        *sql.DB
	currentTx *sql.Tx
}

// NewTransaction returns a new transaction.
func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

// Begin starts a transaction.
func (tx *Transaction) Begin(ctx context.Context) error {
	if tx.currentTx == nil {
		x, err := tx.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin tx: %w", err)
		}
		tx.currentTx = x
	}
	return nil
}

// Commit commits a transaction.
func (tx *Transaction) Commit(ctx context.Context) error {
	if tx.currentTx != nil {
		if err := tx.currentTx.Commit(); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}
		tx.currentTx = nil
	}
	return nil
}

// Rollback rollbacks a transaction.
func (tx *Transaction) Rollback(ctx context.Context) error {
	if tx.currentTx != nil {
		if err := tx.currentTx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback: %w", err)
		}
		tx.currentTx = nil
	}
	return nil
}

func (tx *Transaction) getDB() DB {
	if tx.currentTx != nil {
		return tx.currentTx
	}
	return tx.db
}
