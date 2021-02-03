package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

// StatsRepository is a database repository for statistics.
type StatsRepository struct {
	transaction *Transaction
}

// NewStatsRepository returns a new StatsRepository.
func NewStatsRepository(transaction *Transaction) *StatsRepository {
	return &StatsRepository{
		transaction: transaction,
	}
}

// FindExpensesByMonthInYear finds expenses group by month for given year.
func (repo *StatsRepository) FindExpensesByMonthInYear(ctx context.Context, year int) ([]*domain.AmountByYearMonth, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT strftime('%m',recorded_on) m, abs(sum(amount)) a
FROM money_forward_records
WHERE recorded_on BETWEEN ? AND ?
AND amount < 0
GROUP BY m
ORDER BY m`
	findArgs := []interface{}{
		fmt.Sprintf("%04d-01-01", year),
		fmt.Sprintf("%04d-12-31", year),
	}

	rows, err := db.QueryContext(ctx, findQuery, findArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	out := []*domain.AmountByYearMonth{}
	for rows.Next() {
		var monthStr string
		var amount int
		if err := rows.Scan(&monthStr, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert month string to int: %w", err)
		}
		out = append(out, &domain.AmountByYearMonth{
			Year:   year,
			Month:  time.Month(month),
			Amount: amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return out, nil
}

// FindExpensesByMonthAndCategoryInYear finds expenses group by month and category for given year.
func (repo *StatsRepository) FindExpensesByMonthAndCategoryInYear(ctx context.Context, year int) ([]*domain.AmountByYearMonthCategory, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT strftime('%m',recorded_on) m, c2.id, abs(sum(amount)) a
FROM money_forward_records m
INNER JOIN categories c1 on (m.category_id = c1.id)
INNER JOIN categories c2 on (c1.parent_id = c2.id)
WHERE recorded_on BETWEEN ? AND ?
AND amount < 0
GROUP BY m, c2.id
ORDER BY m, c2.parent_id, c2.display_order, c2.id`
	findArgs := []interface{}{
		fmt.Sprintf("%04d-01-01", year),
		fmt.Sprintf("%04d-12-31", year),
	}

	rows, err := db.QueryContext(ctx, findQuery, findArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	out := []*domain.AmountByYearMonthCategory{}
	for rows.Next() {
		var monthStr string
		var amount int
		var categoryID string
		if err := rows.Scan(&monthStr, &categoryID, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert month string to int: %w", err)
		}
		out = append(out, &domain.AmountByYearMonthCategory{
			Year:       year,
			Month:      time.Month(month),
			CategoryID: categoryID,
			Amount:     amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return out, nil
}

// FindExpensesByCategoryInYearAndMonth
