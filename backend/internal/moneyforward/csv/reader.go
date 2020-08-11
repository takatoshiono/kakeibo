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

type Reader struct {
	reader *csv.Reader

	// If IgnoreHeader is true, ignore header record when read. Default is true.
	IgnoreHeader bool

	// If isHeaderIgnored is true, the header has already ignored.
	isHeaderIgnored bool
}

func NewReader(r io.Reader) *Reader {
	rr := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	return &Reader{
		reader:       csv.NewReader(rr),
		IgnoreHeader: true,
	}
}

func (r *Reader) Read() (record []string, err error) {
	fields, err := r.reader.Read()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	if r.IgnoreHeader && !r.isHeaderIgnored {
		if fields[0] == headerLabels[0] {
			r.isHeaderIgnored = true
			return r.Read()
		}
	}
	return fields, nil
}

func (r *Reader) ReadAll() (records [][]string, err error) {
	return r.reader.ReadAll()
}
