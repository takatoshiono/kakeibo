package database

import (
	"context"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

// MasterRepository is a database repository for master data.
type MasterRepository struct {
	transaction *Transaction
}

// NewMasterRepository returns a new MasterRepository.
func NewMasterRepository(transaction *Transaction) *MasterRepository {
	return &MasterRepository{
		transaction: transaction,
	}
}

// CreateOrFindSource creates or finds a source.
func (repo *MasterRepository) CreateOrFindSource(ctx context.Context, name string) (*domain.Source, error) {
	// TODO: implement
	return nil, nil
}

// CreateOrFindCategory creates or finds a category.
func (repo *MasterRepository) CreateOrFindCategory(ctx context.Context, name string, level domain.CategoryLevel, parentID string) (*domain.Category, error) {
	// TODO: implement
	return nil, nil
}
