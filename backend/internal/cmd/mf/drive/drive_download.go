package drive

import (
	"errors"

	"github.com/spf13/cobra"
)

// NewCmdDriveDownload creates the `drive download` command.
func NewCmdDriveDownload(o *DownloadOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a file from a folder in Google Drive",
		Long: `This command download a file from a folder in Google Drive. For example:

	mf drive download --file FILE_PATH --parent GOOGLE_DRIVE_FOLDER_ID

This command uses Google Application Default Credentials for authentication.
So please to set GOOGLE_APPLICATION_CREDENTIALS environment variable.
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	return cmd
}

// DownloadOption is the option for the `drive download` command.
type DownloadOption struct {
}

// Validate checks options.
func (o *DownloadOption) Validate() error {
	return nil
}

// Run executes the `drive download` command.
func (o *DownloadOption) Run() error {
	return errors.New("not implemented")
}
