package mf

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/takatoshiono/kakeibo/backend/internal/moneyforward"
)

var csvCmd = &cobra.Command{
	Use:   "csv <command>",
	Short: "Download csv files from Money Forward ME",
	Long:  `Work with csv of Money Forward ME`,
}

var csvDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download csv files from Money Forward ME",
	Long:  `This command download files frOm Money Forward ME`,
	RunE:  csvDownload,
}

func csvDownload(cmd *cobra.Command, args []string) error {
	email := os.Getenv("MF_EMAIL")
	password := os.Getenv("MF_PASSWORD")

	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		return fmt.Errorf("failed to get year: %w", err)
	}
	month, err := cmd.Flags().GetInt("month")
	if err != nil {
		return fmt.Errorf("failed to get month: %w", err)
	}
	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		return fmt.Errorf("failed to get filename: %w", err)
	}

	mf := moneyforward.New(email, password)
	ctx := context.Background()
	if err := mf.Login(ctx); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	r, err := mf.DownloadCSV(ctx, year, time.Month(month))
	if err != nil {
		return fmt.Errorf("failed to download csv: %w", err)
	}
	defer r.Close()

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	return nil
}

func init() {
	now := time.Now()
	csvCmd.AddCommand(csvDownloadCmd)
	csvDownloadCmd.Flags().IntP("year", "y", now.Year(), "year. format is YYYY")
	csvDownloadCmd.Flags().IntP("month", "m", int(now.Month()), "month. format is 1 to 12")
	csvDownloadCmd.Flags().StringP("filename", "f", "out.csv", "output filename")
}
