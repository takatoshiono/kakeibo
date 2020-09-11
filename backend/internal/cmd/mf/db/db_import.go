package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // support `sqlite3` driver
	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/moneyforward/csv"
	"github.com/takatoshiono/kakeibo/backend/internal/repository/database"
	"github.com/takatoshiono/kakeibo/backend/internal/usecase"
)

// NewCmdDBImport creates the `db import` command.
func NewCmdDBImport() *cobra.Command {
	o := &DBImportOptions{
		DriverName: os.Getenv("DB_DRIVER_NAME"),
		DSN:        os.Getenv("DB_DSN"),
	}

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import files to database",
		Long:  `This command import files to database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().StringVarP(&o.FileName, "file", "f", "in.csv", "input file name")

	return cmd
}

// DBImportOptions creates new DBImportOptions for the `db import` command.
type DBImportOptions struct {
	DriverName string
	DSN        string
	FileName   string
}

// Run executes the `db import` command.
func (o *DBImportOptions) Run() error {
	f, err := os.Open(o.FileName)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	db, err := sql.Open(o.DriverName, o.DSN)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer db.Close()
	transaction := database.NewTransaction(db)
	masterRepo := database.NewMasterRepository(transaction)
	mfRepo := database.NewMoneyForwardRepository(transaction)

	usecase := usecase.NewImportMoneyForwardRecords(reader, transaction, masterRepo, mfRepo)
	ctx := context.Background()
	if err := usecase.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return nil
}
