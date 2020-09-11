package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

// MoneyForwardRepository is a database repository for money forward data.
type MoneyForwardRepository struct {
	transaction *Transaction
}

// NewMoneyForwardRepository returns a new MoneyForwardRepository.
func NewMoneyForwardRepository(transaction *Transaction) *MoneyForwardRepository {
	return &MoneyForwardRepository{
		transaction: transaction,
	}
}

// CreateOrUpdateRecord creates or updates a money forward record.
func (repo *MoneyForwardRepository) CreateOrUpdateRecord(ctx context.Context, record *domain.MoneyForwardRecord) error {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id FROM money_forward_records WHERE id = ?`
	findArgs := []interface{}{record.ID}

	var id string
	err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		if err := repo.createRecord(ctx, record); err != nil {
			return fmt.Errorf("failed to create record: %w", err)
		}
		return nil
	case err != nil:
		return fmt.Errorf("failed to scan: %w", err)
	default:
		if err := repo.updateRecord(ctx, record); err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}
		return nil
	}
}

func (repo *MoneyForwardRepository) createRecord(ctx context.Context, record *domain.MoneyForwardRecord) error {
	db := repo.transaction.getDB()

	now := time.Now()

	const insertQuery = `
INSERT INTO money_forward_records(id, recorded_on, title, amount, source_id, category_id, memo, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	insertArgs := []interface{}{record.ID, record.RecordedOn.String(), record.Title, record.Amount, record.SourceID, record.Category2ID, record.Memo, now, now}

	if _, err := db.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (repo *MoneyForwardRepository) updateRecord(ctx context.Context, record *domain.MoneyForwardRecord) error {
	db := repo.transaction.getDB()

	now := time.Now()

	const updateQuery = `
UPDATE money_forward_records
SET recorded_on = ?, title = ?, amount = ?, source_id = ?, category_id = ?, memo = ?, updated_at = ?
WHERE id = ?`
	updateArgs := []interface{}{record.RecordedOn.String(), record.Title, record.Amount, record.SourceID, record.Category2ID, record.Memo, now, record.ID}

	if _, err := db.ExecContext(ctx, updateQuery, updateArgs...); err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}
