package csv

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/moneyforward"
)

// NewCmdCSVDownload creates the `csv download` command.
func NewCmdCSVDownload(o *DownloadOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download csv files from Money Forward ME",
		Long: `This command download files from Money Forward ME.

If from and to options are specified, download csv files for period "from" to "to".
The output file name has "-YYYYMM" suffix.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Parse(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	now := time.Now()
	cmd.Flags().IntVarP(&o.year, "year", "y", now.Year(), "year. format is YYYY")
	cmd.Flags().IntVarP(&o.month, "month", "m", int(now.Month()), "month. format is 1 to 12")
	cmd.Flags().StringVar(&o.from, "from", "", "from year and month. format is YYYYMM")
	cmd.Flags().StringVar(&o.to, "to", "", "to year and month. format is YYYYMM")
	cmd.Flags().StringVarP(&o.fileName, "filename", "f", "mf.csv", "output filename. if from and to are specified it will convert as mf-YYYYMM.csv")

	return cmd
}

// DownloadOption is the option for the `csv download` command.
type DownloadOption struct {
	MoneyForwardEmail    string
	MoneyForwardPassword string
	year                 int
	month                int
	from                 string
	to                   string
	fileName             string

	fromTime time.Time
	toTime   time.Time
}

// Validate checks options.
func (o *DownloadOption) Validate() error {
	if o.year == 0 {
		return fmt.Errorf("year must not be 0")
	}
	if o.month < 1 || 12 < o.month {
		return fmt.Errorf("month must be between 1 and 12")
	}
	if o.from == "" && o.to != "" || o.from != "" && o.to == "" {
		return fmt.Errorf("both from and to must be set")
	}
	if o.from != "" && o.to != "" {
		const yearMonthFormat = "200601"
		f, err := time.Parse(yearMonthFormat, o.from)
		if err != nil {
			return fmt.Errorf("from must be YYYYMM format")
		}
		t, err := time.Parse(yearMonthFormat, o.to)
		if err != nil {
			return fmt.Errorf("to must be YYYYMM format")
		}
		if t.Before(f) {
			return fmt.Errorf("to must be equal or after from")
		}
		o.fromTime, o.toTime = f, t
	}
	if o.fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}
	return nil
}

// Parse parses options.
func (o *DownloadOption) Parse() error {
	if o.fromTime.IsZero() && o.toTime.IsZero() {
		o.fromTime = time.Date(o.year, time.Month(o.month), 1, 0, 0, 0, 0, time.UTC)
		o.toTime = o.fromTime
	}
	return nil
}

// Run executes the `csv download` command.
func (o *DownloadOption) Run() error {
	mf := moneyforward.New(o.MoneyForwardEmail, o.MoneyForwardPassword)
	ctx := context.Background()
	// TODO: add timeout to context

	log.Printf("login as %s\n", o.MoneyForwardEmail)

	if err := mf.Login(ctx); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	tm := o.fromTime
	for !tm.After(o.toTime) {
		log.Printf("download csv for %04d%02d\n", tm.Year(), tm.Month())
		r, err := mf.DownloadCSV(ctx, tm.Year(), tm.Month())
		if err != nil {
			return fmt.Errorf("failed to download csv: %w", err)
		}
		defer r.Close()

		fileName := o.convertFileName(tm)
		log.Printf("create %s\n", fileName)
		f, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("failed to open: %w", err)
		}
		defer f.Close()
		if _, err := io.Copy(f, r); err != nil {
			return fmt.Errorf("failed to copy: %w", err)
		}

		tm = tm.AddDate(0, 1, 0)
	}

	return nil
}

func (o *DownloadOption) convertFileName(t time.Time) string {
	dir, file := filepath.Split(o.fileName)
	if strings.Index(file, ".") != -1 {
		sub := strings.Split(file, ".")
		return filepath.Join(dir, fmt.Sprintf("%s-%04d%02d.%s", sub[0], t.Year(), t.Month(), sub[1]))
	}
	return filepath.Join(dir, fmt.Sprintf("%s-%04d%02d", file, t.Year(), t.Month()))
}
