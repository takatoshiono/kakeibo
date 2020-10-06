package drive

import (
	"github.com/spf13/cobra"
)

type (
	// Options is the collection of options for the `db` command and its sub command.
	Options struct {
		UploadOption   *UploadOption
		DownloadOption *DownloadOption
	}
)

// NewCmdDrive creates the `drive` command.
func NewCmdDrive(o *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drive <command>",
		Short: "Upload and download files to Google Drive",
		Long:  `Work with Google Drive`,
	}

	cmd.AddCommand(NewCmdDriveUpload(o.UploadOption))
	cmd.AddCommand(NewCmdDriveDownload(o.DownloadOption))

	return cmd
}
