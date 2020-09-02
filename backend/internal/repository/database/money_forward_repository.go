package database

import (
	"context"

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
	// TODO: implement
	return nil
}
