package mf

import "context"

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// DownloadAndUploadCSV consumes a Pub/Sub message and download CSV from MoneyForward ME
// then upload it to the Google Drive folder.
func DownloadAndUploadCSV(ctx context.Context, m PubSubMessage) error {
	return nil
}
