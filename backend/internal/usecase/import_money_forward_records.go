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

		if !record.IsRecordToSave() {
			continue
		}

		if err := u.transaction.Begin(ctx); err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		source, err := u.masterRepo.FindOrCreateSource(ctx, record.SourceName)
		if err != nil {
			wrapErr := fmt.Errorf("failed to find or create source: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		category1, err := u.masterRepo.FindOrCreateCategory(ctx, record.Category1Name, domain.CategoryLevel1, "")
		if err != nil {
			wrapErr := fmt.Errorf("failed to find or create category level 1: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		category2, err := u.masterRepo.FindOrCreateCategory(ctx, record.Category2Name, domain.CategoryLevel2, category1.ID)
		if err != nil {
			wrapErr := fmt.Errorf("failed to find or create category level 2: %w", err)
			if err := u.transaction.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback: %w", wrapErr)
			}
			return wrapErr
		}

		record.SourceID = source.ID
		record.Category1ID = category1.ID
		record.Category2ID = category2.ID

		if err := u.mfRepo.CreateOrUpdateRecord(ctx, record); err != nil {
			wrapErr := fmt.Errorf("failed to create or update record: %w", err)
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
