package mf

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/takatoshiono/kakeibo/backend/pkg/googledrive"
)

var driveCmd = &cobra.Command{
	Use:   "drive <command>",
	Short: "Upload and download files to Google Drive",
	Long:  `Work with Google Drive`,
}

var driveUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to a folder in Google Drive",
	Long: `This command upload a file to a folder in Google Drive. For example:

	mf drive upload --file FILE_PATH --parent GOOGLE_DRIVE_FOLDER_ID
	
This command uses Google Application Default Credentials for authentication.
So please to set GOOGLE_APPLICATION_CREDENTIALS environment variable.
	`,
	RunE: driveUpload,
}

var driveDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from a folder in Google Drive",
	Long: `This command download a file from a folder in Google Drive. For example:

	mf drive download --file FILE_PATH --parent GOOGLE_DRIVE_FOLDER_ID
	
This command uses Google Application Default Credentials for authentication.
So please to set GOOGLE_APPLICATION_CREDENTIALS environment variable.
	`,
	RunE: driveDownload,
}

func driveUpload(cmd *cobra.Command, args []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	parentID, err := cmd.Flags().GetString("parent")
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open: %v", err)
	}
	defer f.Close()

	ctx := context.Background()
	d, err := googledrive.New(ctx)
	if err != nil {
		log.Fatalf("failed to new googledrive: %v", err)
	}

	fileID, err := d.CreateFile(ctx, f, filepath.Base(filePath), "text/plain", parentID)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	fmt.Printf("created %s\n", fileID)

	return nil
}

func driveDownload(cmd *cobra.Command, args []string) error {
	return errors.New("not implemented")
}

func init() {
	driveCmd.AddCommand(driveUploadCmd)
	driveUploadCmd.Flags().StringP("file", "f", "", "upload file path")
	driveUploadCmd.Flags().StringP("parent", "p", "", "parent folder id")
	driveUploadCmd.MarkFlagRequired("file")
	driveUploadCmd.MarkFlagRequired("parent")

	driveCmd.AddCommand(driveDownloadCmd)
}
