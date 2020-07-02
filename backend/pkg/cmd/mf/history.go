package mf

import (
	"errors"

	"github.com/spf13/cobra"
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
	return errors.New("not implemented")
}

func init() {
	historyCmd.AddCommand(historyDownloadCmd)
}
