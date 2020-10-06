package drive

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/takatoshiono/kakeibo/backend/internal/googledrive"
)

// NewCmdDriveUpload creates the `drive upload` command.
func NewCmdDriveUpload(o *UploadOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to a folder in Google Drive",
		Long: `This command upload a file to a folder in Google Drive. For example:

	mf drive upload --file FILE_PATH --parent GOOGLE_DRIVE_FOLDER_ID

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

	cmd.Flags().StringP("file", "f", "", "upload file path")
	cmd.Flags().StringP("parent", "p", "", "parent folder id")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("parent")

	return cmd
}

// UploadOption is the option for the `drive upload` command.
type UploadOption struct {
	fileName string
	parentID string
}

// Validate checks options.
func (o *UploadOption) Validate() error {
	if o.fileName == "" {
		return fmt.Errorf("file name must not be empty")
	}
	if o.parentID == "" {
		return fmt.Errorf("parent id must not be empty")
	}
	return nil
}

// Run executes the `drive upload` command.
func (o *UploadOption) Run() error {
	f, err := os.Open(o.fileName)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer f.Close()

	ctx := context.Background()
	d, err := googledrive.New(ctx)
	if err != nil {
		log.Fatalf("failed to new googledrive: %v", err)
	}

	fileID, err := d.CreateFile(ctx, f, filepath.Base(o.fileName), "text/plain", o.parentID)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	fmt.Printf("created %s\n", fileID)

	return nil
}
