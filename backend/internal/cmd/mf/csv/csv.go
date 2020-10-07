package csv

import (
	"github.com/spf13/cobra"
)

type (
	// Options is the collection of options for the `db` command and its sub command.
	Options struct {
		DownloadOption *DownloadOption
	}
)

// NewCmdCSV creates the `csv` command.
func NewCmdCSV(o *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "csv <command>",
		Short: "Download csv files from Money Forward ME",
		Long:  `Work with csv of Money Forward ME`,
	}

	cmd.AddCommand(NewCmdCSVDownload(o.DownloadOption))

	return cmd
}
