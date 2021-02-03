package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// StatsMoneyForwardRecordsArgs holds arguments to execute StatsMoneyForwardRecords.
type StatsMoneyForwardRecordsArgs struct {
	Year int
}

// StatsMoneyForwardRecords is a usecase for money forward records.
type StatsMoneyForwardRecords struct {
	transaction Transaction
	statsRepo   StatsRepository
	w           io.Writer
}

// NewStatsMoneyForwardRecords returns a new StatsMoneyForwardRecords usecase.
func NewStatsMoneyForwardRecords(transaction Transaction, statsRepo StatsRepository, w io.Writer) *StatsMoneyForwardRecords {
	return &StatsMoneyForwardRecords{
		transaction: transaction,
		statsRepo:   statsRepo,
		w:           w,
	}
}

// Execute executes the usecase.
func (u *StatsMoneyForwardRecords) Execute(ctx context.Context, queryName string, args *StatsMoneyForwardRecordsArgs) error {
	out := [][]string{}
	switch queryName {
	case "ExpensesByMonth":
		res, err := u.statsRepo.FindExpensesByMonthInYear(ctx, args.Year)
		if err != nil {
			return fmt.Errorf("failed to find amount expended by month: %w", err)
		}
		for _, r := range res {
			out = append(out, []string{
				strconv.Itoa(int(r.Month)),
				strconv.Itoa(r.Amount),
			})
		}
		if err := u.outputCSV(out); err != nil {
			return fmt.Errorf("failed to output csv: %w", err)
		}
	default:
		return fmt.Errorf("unknown query name '%v'", queryName)
	}
	return nil
}

func (u *StatsMoneyForwardRecords) outputCSV(records [][]string) error {
	w := csv.NewWriter(u.w)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed to flush: %w", err)
	}
	return nil
}
