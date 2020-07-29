package googledrive

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type GoogleDrive struct {
	service *drive.Service
}

func New(ctx context.Context) (*GoogleDrive, error) {
	service, err := drive.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to new service: %w", err)
	}
	return &GoogleDrive{
		service: service,
	}, nil
}

func (d *GoogleDrive) CreateFile(ctx context.Context, r io.Reader, filename string, ctype string, parentID string) (string, error) {
	file := &drive.File{
		Name:    filename,
		Parents: []string{parentID},
	}
	c := d.service.Files.Create(file)
	c.Context(ctx)
	f, err := c.Media(r, googleapi.ContentType(ctype)).Do()
	if err != nil {
		return "", fmt.Errorf("failed to do files create call: %w", err)
	}
	return f.Id, nil
}
