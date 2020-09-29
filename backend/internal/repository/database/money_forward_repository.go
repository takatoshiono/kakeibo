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

// FindRecord finds the record by id.
func (repo *MoneyForwardRepository) FindRecord(ctx context.Context, id string) (*domain.MoneyForwardRecord, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, recorded_on, title, amount, source_id, category_id, memo
FROM money_forward_records WHERE id = ?`
	findArgs := []interface{}{id}

	r := &domain.MoneyForwardRecord{}
	var recordedOn string
	if err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&r.ID, &recordedOn, &r.Title, &r.Amount, &r.SourceID, &r.Category2ID, &r.Memo); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	t, err := time.Parse(time.RFC3339, recordedOn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recorded_on: %w", err)
	}
	r.RecordedOn = t

	return r, nil
}

func (repo *MoneyForwardRepository) createRecord(ctx context.Context, record *domain.MoneyForwardRecord) error {
	db := repo.transaction.getDB()

	now := time.Now()

	const insertQuery = `
INSERT INTO money_forward_records(id, recorded_on, title, amount, source_id, category_id, memo, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// TODO: RecordedOnは型を定義してRecordedOn.String() とかで YYYY-MM-DD を取得できるとよさそう
	insertArgs := []interface{}{record.ID, record.RecordedOn.Format("2006-01-02"), record.Title, record.Amount, record.SourceID, record.Category2ID, record.Memo, now, now}

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
	updateArgs := []interface{}{record.RecordedOn.Format("2006-01-02"), record.Title, record.Amount, record.SourceID, record.Category2ID, record.Memo, now, record.ID}

	if _, err := db.ExecContext(ctx, updateQuery, updateArgs...); err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}
