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

// FindOrCreateSource creates or finds a source.
func (repo *MasterRepository) FindOrCreateSource(ctx context.Context, name string) (*domain.Source, error) {
	// TODO: implement
	return &domain.Source{}, nil
}

// FindOrCreateCategory createFindOrinds a category.
func (repo *MasterRepository) FindOrCreateCategory(ctx context.Context, name string, level domain.CategoryLevel, parentID string) (*domain.Category, error) {
	// TODO: implement
	return &domain.Category{}, nil
}
