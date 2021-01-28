package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // support `sqlite3` driver
	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/repository/database"
	"github.com/takatoshiono/kakeibo/backend/internal/usecase"
)

// NewCmdDBStats creates the `db stats` command.
func NewCmdDBStats(o *StatsOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Get statistics",
		Long:  `This command get statistics`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.queryName, "query", "q", "", "stats query name")
	cmd.Flags().IntVarP(&o.year, "year", "y", 0, "year")

	return cmd
}

// StatsOption is the option for the `db stats` command.
type StatsOption struct {
	DriverName string
	DSN        string
	queryName  string
	year       int
}

// Validate checks options.
func (o *StatsOption) Validate() error {
	if o.queryName == "" {
		return fmt.Errorf("query name must not be empty")
	}
	return nil
}

// Run executes the `db stats` command.
func (o *StatsOption) Run() error {
	db, err := sql.Open(o.DriverName, o.DSN)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer db.Close()
	transaction := database.NewTransaction(db)
	statsRepo := database.NewStatsRepository(transaction)

	u := usecase.NewStatsMoneyForwardRecords(transaction, statsRepo, os.Stdout)
	args := &usecase.StatsMoneyForwardRecordsArgs{
		Year: o.year,
	}
	ctx := context.Background()
	if err := u.Execute(ctx, o.queryName, args); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return nil
}
