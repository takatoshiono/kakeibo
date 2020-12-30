package db

import (
	"github.com/spf13/cobra"
)

type (
	// Options is the collection of options for the `db` command and its sub command.
	Options struct {
		ImportOption *ImportOption
		StatsOption  *StatsOption
	}
)

// NewCmdDB creates the `db` command.
func NewCmdDB(o *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db <command>",
		Short: "Import files to database",
		Long:  `Work with database`,
	}

	// TODO: Add delete command
	cmd.AddCommand(NewCmdDBImport(o.ImportOption))
	cmd.AddCommand(NewCmdDBStats(o.StatsOption))

	return cmd
}
