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

// FindExpensesByMonthInYear finds amount expended group by month for given year.
func (repo *StatsRepository) FindExpensesByMonthInYear(ctx context.Context, year int) ([]*domain.AmountExpendedByMonth, error) {
	db := repo.transaction.getDB()

	const findQuery = `
SELECT strftime('%m',recorded_on) m, abs(sum(amount)) a
FROM money_forward_records
WHERE recorded_on BETWEEN ? AND ?
AND amount < 0
GROUP BY m`
	findArgs := []interface{}{
		fmt.Sprintf("%04d-01-01", year),
		fmt.Sprintf("%04d-12-31", year),
	}

	rows, err := db.QueryContext(ctx, findQuery, findArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	out := []*domain.AmountExpendedByMonth{}
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
		out = append(out, &domain.AmountExpendedByMonth{
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

// FindExpensesByMonthInYear
// FindExpensesByMonthAndCategoryInYear
// FindExpensesByCategoryInYearAndMonth
