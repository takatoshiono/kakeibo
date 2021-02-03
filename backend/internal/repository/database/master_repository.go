package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

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

// FindOrCreateSource finds or creates a source.
func (repo *MasterRepository) FindOrCreateSource(ctx context.Context, name string) (*domain.Source, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, name, display_order FROM sources WHERE name = ?`
	findArgs := []interface{}{name}

	s := &domain.Source{}
	err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&s.ID, &s.Name, &s.DisplayOrder)
	switch {
	case err == sql.ErrNoRows:
		// pass through
	case err != nil:
		return nil, fmt.Errorf("failed to scan: %w", err)
	default:
		return s, nil
	}

	now := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to get a random uuid: %w", err)
	}
	displayOrder := 0

	const insertQuery = `
INSERT INTO sources(id, name, display_order, created_at, updated_at)
VALUES(?, ?, ?, ?, ?)`
	insertArgs := []interface{}{uuid.String(), name, displayOrder, now, now}

	if _, err := db.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return &domain.Source{
		ID:           uuid.String(),
		Name:         name,
		DisplayOrder: displayOrder,
	}, nil
}

// FindOrCreateCategory finds or creates a category.
func (repo *MasterRepository) FindOrCreateCategory(ctx context.Context, name string, level domain.CategoryLevel, parentID string) (*domain.Category, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, name, level, display_order, parent_id FROM categories WHERE name = ? AND level = ?`
	findArgs := []interface{}{name, level}

	// TODO: database packageのstructを定義してdomain packageのstructに変換したほうが良い？
	c := &domain.Category{}
	err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&c.ID, &c.Name, &c.Level, &c.DisplayOrder, &c.ParentID)
	switch {
	case err == sql.ErrNoRows:
		// pass through
	case err != nil:
		return nil, fmt.Errorf("failed to scan: %w", err)
	default:
		return c, nil
	}

	now := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to get a random uuid: %w", err)
	}
	displayOrder := 0

	const insertQuery = `
INSERT INTO categories(id, name, display_order, level, parent_id, created_at, updated_at)
VALUES(?, ?, ?, ?, ?, ?, ?)`
	insertArgs := []interface{}{uuid.String(), name, displayOrder, level, parentID, now, now}

	if _, err := db.ExecContext(ctx, insertQuery, insertArgs...); err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return &domain.Category{
		ID:           uuid.String(),
		Name:         name,
		DisplayOrder: displayOrder,
		Level:        level,
		ParentID:     parentID,
	}, nil
}

// FindSourceByName finds the source by name.
func (repo *MasterRepository) FindSourceByName(ctx context.Context, name string) (*domain.Source, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, name, display_order FROM sources WHERE name = ?`
	findArgs := []interface{}{name}

	s := &domain.Source{}
	if err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&s.ID, &s.Name, &s.DisplayOrder); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return s, nil
}

// FindCategoryByID finds the category by id.
func (repo *MasterRepository) FindCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, name, level, display_order, parent_id FROM categories WHERE id = ?`
	findArgs := []interface{}{id}

	c := &domain.Category{}
	if err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&c.ID, &c.Name, &c.Level, &c.DisplayOrder, &c.ParentID); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return c, nil
}

// FindCategoryByNameAndLevel finds the category by name and level.
func (repo *MasterRepository) FindCategoryByNameAndLevel(ctx context.Context, name string, level domain.CategoryLevel) (*domain.Category, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT id, name, level, display_order, parent_id FROM categories WHERE name = ? AND level = ?`
	findArgs := []interface{}{name, level}

	c := &domain.Category{}
	if err := db.QueryRowContext(ctx, findQuery, findArgs...).Scan(&c.ID, &c.Name, &c.Level, &c.DisplayOrder, &c.ParentID); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return c, nil
}
