package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/takatoshiono/kakeibo/backend/internal/convert/moneyforward"
	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

// ImportMoneyForwardRecords is a usecase for money forward records.
type ImportMoneyForwardRecords struct {
	reader      MoneyForwardCSVReader
	transaction Transaction
	masterRepo  MasterRepository
	mfRepo      MoneyForwardRepository
}

// NewImportMoneyForwardRecords returns a new ImportMoneyForwardRecords usecase.
func NewImportMoneyForwardRecords(reader MoneyForwardCSVReader, transaction Transaction, masterRepo MasterRepository, mfRepo MoneyForwardRepository) *ImportMoneyForwardRecords {
	return &ImportMoneyForwardRecords{
		reader:      reader,
		transaction: transaction,
		masterRepo:  masterRepo,
		mfRepo:      mfRepo,
	}
}

// Execute executes the usecase.
func (u *ImportMoneyForwardRecords) Execute(ctx context.Context) error {
	for {
		fields, err := u.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}

		record, err := moneyforward.ConvCSVToDomain(fields)
		if err != nil {
			return fmt.Errorf("failed to convert: %w", err)
		}

		fmt.Println(record)

		if err := u.transaction.Begin(ctx); err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := u.masterRepo.CreateOrFindSource(ctx, record.Source); err != nil {
			wrapErr := fmt.Errorf("failed to create or find source: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		category1, err := u.masterRepo.CreateOrFindCategory(ctx, record.Category1, domain.CategoryLevel1, "")
		if err != nil {
			wrapErr := fmt.Errorf("failed to create or find source: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		if _, err := u.masterRepo.CreateOrFindCategory(ctx, record.Category2, domain.CategoryLevel2, category1.ID); err != nil {
			wrapErr := fmt.Errorf("failed to create or find source: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		if err := u.mfRepo.CreateOrUpdateRecord(ctx, record); err != nil {
			wrapErr := fmt.Errorf("failed to CreateOrUpdateRecord record: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		if err := u.transaction.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}
	}
	return nil
}
