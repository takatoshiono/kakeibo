package mf

import (
	"errors"

	"github.com/spf13/cobra"
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
	return errors.New("not implemented")
}

func init() {
	dbCmd.AddCommand(dbImportCmd)
}
