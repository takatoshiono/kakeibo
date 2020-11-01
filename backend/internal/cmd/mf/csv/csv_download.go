package csv

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/moneyforward"
)

// NewCmdCSVDownload creates the `csv download` command.
func NewCmdCSVDownload(o *DownloadOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download csv files from Money Forward ME",
		Long:  `This command download files frOm Money Forward ME`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	now := time.Now()
	cmd.Flags().IntVarP(&o.year, "year", "y", now.Year(), "year. format is YYYY")
	cmd.Flags().IntVarP(&o.month, "month", "m", int(now.Month()), "month. format is 1 to 12")
	cmd.Flags().StringVarP(&o.fileName, "filename", "f", "out.csv", "output filename")

	return cmd
}

// DownloadOption is the option for the `csv download` command.
type DownloadOption struct {
	MoneyForwardEmail    string
	MoneyForwardPassword string
	year                 int
	month                int
	fileName             string
}

// Validate checks options.
func (o *DownloadOption) Validate() error {
	if o.year == 0 {
		return fmt.Errorf("year must not be 0")
	}
	if o.month < 1 || 12 < o.month {
		return fmt.Errorf("month must be between 1 and 12")
	}
	if o.fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}
	return nil
}

// Run executes the `csv download` command.
func (o *DownloadOption) Run() error {
	mf := moneyforward.New(o.MoneyForwardEmail, o.MoneyForwardPassword)
	ctx := context.Background()
	if err := mf.Login(ctx); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	r, err := mf.DownloadCSV(ctx, o.year, time.Month(o.month))
	if err != nil {
		return fmt.Errorf("failed to download csv: %w", err)
	}
	defer r.Close()

	f, err := os.Create(o.fileName)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	return nil
}
