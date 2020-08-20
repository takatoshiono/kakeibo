package csv

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestNewReader(t *testing.T) {
	t.Run("default IgnoreHeader is true", func(t *testing.T) {
		r := NewReader(strings.NewReader(""))
		if got, want := r.IgnoreHeader, true; got != want {
			t.Errorf("want IgnoreHeader %v, but got %v", want, got)
		}
	})
}

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name         string
		r            io.Reader
		ignoreHeader bool
		want         []string
		wantErr      bool
	}{
		{
			name: "ignore header record",
			r: transform.NewReader(
				strings.NewReader(`"計算対象","日付","内容","金額（円）","保有金融機関","大項目","中項目","メモ","振替","ID"
"1","2020/07/25","西友","-3612","セゾンカード","食費","食料品","","0","qi03Xo5JDVYjZC2HqFA9Sg"`),
				japanese.ShiftJIS.NewEncoder()),
			ignoreHeader: true,
			want: []string{
				"1", "2020/07/25", "西友", "-3612", "セゾンカード", "食費", "食料品", "", "0", "qi03Xo5JDVYjZC2HqFA9Sg",
			},
			wantErr: false,
		},
		{
			name: "not ignore header record",
			r: transform.NewReader(
				strings.NewReader(`"1","2020/07/25","西友","-3612","セゾンカード","食費","食料品","","0","qi03Xo5JDVYjZC2HqFA9Sg"`),
				japanese.ShiftJIS.NewEncoder()),
			ignoreHeader: false,
			want: []string{
				"1", "2020/07/25", "西友", "-3612", "セゾンカード", "食費", "食料品", "", "0", "qi03Xo5JDVYjZC2HqFA9Sg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.r)
			r.IgnoreHeader = tt.ignoreHeader
			got, err := r.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadAll(t *testing.T) {
	tests := []struct {
		name         string
		r            io.Reader
		ignoreHeader bool
		want         [][]string
		wantErr      bool
	}{
		{
			name: "ignore header record",
			r: transform.NewReader(
				strings.NewReader(`"計算対象","日付","内容","金額（円）","保有金融機関","大項目","中項目","メモ","振替","ID"
"1","2020/07/25","西友","-3612","セゾンカード","食費","食料品","","0","qi03Xo5JDVYjZC2HqFA9Sg"
"1","2020/07/22","給料","150000","三菱UFJ銀行","収入","給与","","0","E_nfxaL2_WJLkpzPUUcRbw"`),
				japanese.ShiftJIS.NewEncoder()),
			ignoreHeader: true,
			want: [][]string{
				{"1", "2020/07/25", "西友", "-3612", "セゾンカード", "食費", "食料品", "", "0", "qi03Xo5JDVYjZC2HqFA9Sg"},
				{"1", "2020/07/22", "給料", "150000", "三菱UFJ銀行", "収入", "給与", "", "0", "E_nfxaL2_WJLkpzPUUcRbw"},
			},
			wantErr: false,
		},
		{
			name: "not ignore header record",
			r: transform.NewReader(
				strings.NewReader(`"1","2020/07/25","西友","-3612","セゾンカード","食費","食料品","","0","qi03Xo5JDVYjZC2HqFA9Sg"
"1","2020/07/22","給料","150000","三菱UFJ銀行","収入","給与","","0","E_nfxaL2_WJLkpzPUUcRbw"`),
				japanese.ShiftJIS.NewEncoder()),
			ignoreHeader: false,
			want: [][]string{
				{"1", "2020/07/25", "西友", "-3612", "セゾンカード", "食費", "食料品", "", "0", "qi03Xo5JDVYjZC2HqFA9Sg"},
				{"1", "2020/07/22", "給料", "150000", "三菱UFJ銀行", "収入", "給与", "", "0", "E_nfxaL2_WJLkpzPUUcRbw"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.r)
			r.IgnoreHeader = tt.ignoreHeader
			got, err := r.ReadAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
