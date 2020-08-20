package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var headerLabels = [10]string{
	"計算対象", "日付", "内容", "金額（円）", "保有金融機関", "大項目", "中項目", "メモ", "振替", "ID",
}

// Reader reads records from money forward csv file.
type Reader struct {
	reader *csv.Reader

	// If IgnoreHeader is true, ignore header record when read. Default is true.
	IgnoreHeader bool

	// If isHeaderIgnored is true, the header has already ignored.
	isHeaderIgnored bool
}

// NewReader returns a new reader that reads from r.
func NewReader(r io.Reader) *Reader {
	rr := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	return &Reader{
		reader:       csv.NewReader(rr),
		IgnoreHeader: true,
	}
}

// Read reads one record from r.
func (r *Reader) Read() ([]string, error) {
	record, err := r.reader.Read()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	if r.IgnoreHeader && !r.isHeaderIgnored {
		if r.isHeaderRecord(record) {
			r.isHeaderIgnored = true
			return r.Read()
		}
	}
	return record, nil
}

// ReadAll reads all the remaining records from r.
func (r *Reader) ReadAll() ([][]string, error) {
	records, err := r.reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to ReadAll: %w", err)
	}

	if r.IgnoreHeader && !r.isHeaderIgnored {
		if r.isHeaderRecord(records[0]) {
			r.isHeaderIgnored = true
			return records[1:], nil
		}
	}
	return records, nil
}

func (r *Reader) isHeaderRecord(record []string) bool {
	return record[0] == headerLabels[0]
}
