package mf

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/moneyforward/csv"
	"github.com/takatoshiono/kakeibo/backend/internal/repository/database"
	"github.com/takatoshiono/kakeibo/backend/internal/usecase"
)

var dbCmd = &cobra.Command{
	Use:   "db <command>",
	Short: "Import files to database",
	Long:  `Work with database`,
}

var dbImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import files to database",
	Long:  `This command import files to database`,
	RunE:  dbImport,
}

func dbImport(cmd *cobra.Command, args []string) error {
	driverName := os.Getenv("DB_DRIVER_NAME")
	dsn := os.Getenv("DB_DSN")

	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return fmt.Errorf("failed to get filename: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	db, err := sql.Open(driverName, dsn)
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

func init() {
	dbCmd.AddCommand(dbImportCmd)
	dbImportCmd.Flags().StringP("filename", "f", "in.csv", "input filename")
}
