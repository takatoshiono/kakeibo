package db

import (
	"github.com/spf13/cobra"
)

// NewCmdDB creates the `db` command.
func NewCmdDB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db <command>",
		Short: "Import files to database",
		Long:  `Work with database`,
	}

	cmd.AddCommand(NewCmdDBImport())

	return cmd
}
