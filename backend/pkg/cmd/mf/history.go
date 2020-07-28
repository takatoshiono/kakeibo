package mf

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/takatoshiono/kakeibo/backend/pkg/moneyforward"
)

var historyCmd = &cobra.Command{
	Use:   "history <command>",
	Short: "Download history files from Money Forward ME",
	Long:  `Work with history of Money Forward ME`,
}

var historyDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download history files from Money Forward ME",
	Long:  `This command download files from Money Forward ME`,
	RunE:  historyDownload,
}

func historyDownload(cmd *cobra.Command, args []string) error {
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

	// TODO: validate input arguments

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
	historyCmd.AddCommand(historyDownloadCmd)
	historyDownloadCmd.Flags().IntP("year", "y", 0, "year. format is YYYY")
	historyDownloadCmd.Flags().IntP("month", "m", 0, "month. format is 1 to 12")
	historyDownloadCmd.Flags().StringP("filename", "f", "out.csv", "output filename. default: out.csv")
}
