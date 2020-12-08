package drive

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/internal/googledrive"
)

// NewCmdDriveDelete creates the `drive delete` command.
func NewCmdDriveDelete(o *DeleteOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete files in the Google Drive folder",
		Long: `This command delete files in the Google Drive folder. For example:

	mf drive delete --file FILE_NAME --parent GOOGLE_DRIVE_FOLDER_ID

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

	cmd.Flags().StringVarP(&o.fileName, "file", "f", "", "upload file path")
	cmd.Flags().StringVarP(&o.parentID, "parent", "p", "", "parent folder id")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("parent")

	return cmd
}

// DeleteOption is the option for the `drive delete` command.
type DeleteOption struct {
	fileName string
	parentID string
}

// Validate checks options.
func (o *DeleteOption) Validate() error {
	if o.fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}
	if o.parentID == "" {
		return fmt.Errorf("parent id must not be empty")
	}
	return nil
}

// Run executes the `drive delete` command.
func (o *DeleteOption) Run() error {
	ctx := context.Background()
	d, err := googledrive.New(ctx)
	if err != nil {
		log.Fatalf("failed to new googledrive: %v", err)
	}

	if err := d.DeleteFileByName(ctx, o.fileName, o.parentID); err != nil {
		log.Fatalf("failed to delete file: %v", err)
	}
	fmt.Printf("deleted %s\n", o.fileName)

	return nil
}
