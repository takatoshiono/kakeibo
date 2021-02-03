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
	masterRepo  MasterRepository
	w           io.Writer
}

// NewStatsMoneyForwardRecords returns a new StatsMoneyForwardRecords usecase.
func NewStatsMoneyForwardRecords(transaction Transaction, statsRepo StatsRepository, masterRepo MasterRepository, w io.Writer) *StatsMoneyForwardRecords {
	return &StatsMoneyForwardRecords{
		transaction: transaction,
		statsRepo:   statsRepo,
		masterRepo:  masterRepo,
		w:           w,
	}
}

// Execute executes the usecase.
func (u *StatsMoneyForwardRecords) Execute(ctx context.Context, queryName string, args *StatsMoneyForwardRecordsArgs) error {
	switch queryName {
	case "ExpensesByMonth":
		return u.executeExpensesByMonth(ctx, args.Year)
	case "ExpensesByMonthAndCategory":
		return u.executeExpensesByMonthAndCategory(ctx, args.Year)
	default:
		return fmt.Errorf("unknown query name '%v'", queryName)
	}
}

func (u *StatsMoneyForwardRecords) executeExpensesByMonth(ctx context.Context, year int) error {
	res, err := u.statsRepo.FindExpensesByMonthInYear(ctx, year)
	if err != nil {
		return fmt.Errorf("failed to find expenses by month: %w", err)
	}

	out := [][]string{}
	for _, r := range res {
		out = append(out, []string{
			strconv.Itoa(int(r.Month)),
			strconv.Itoa(r.Amount),
		})
	}

	if err := u.outputCSV(out); err != nil {
		return fmt.Errorf("failed to output csv: %w", err)
	}

	return nil
}

func (u *StatsMoneyForwardRecords) executeExpensesByMonthAndCategory(ctx context.Context, year int) error {
	res, err := u.statsRepo.FindExpensesByMonthAndCategoryInYear(ctx, year)
	if err != nil {
		return fmt.Errorf("failed to find expenses by month and category: %w", err)
	}

	out := [][]string{}
	for _, r := range res {
		c, err := u.masterRepo.FindCategoryByID(ctx, r.CategoryID)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}
		out = append(out, []string{
			strconv.Itoa(int(r.Month)),
			c.Name,
			strconv.Itoa(r.Amount),
		})
	}

	if err := u.outputCSV(out); err != nil {
		return fmt.Errorf("failed to output csv: %w", err)
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
