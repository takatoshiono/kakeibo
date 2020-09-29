package moneyforward

import (
	"reflect"
	"testing"
	"time"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

func TestConvCSVToDomain(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.MoneyForwardRecord
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				fields: []string{
					"1", "2020/07/25", "西友", "-3612", "セゾンカード", "食費", "食料品", "", "0", "qi03Xo5JDVYjZC2HqFA9Sg",
				},
			},
			want: &domain.MoneyForwardRecord{
				ID:                  "qi03Xo5JDVYjZC2HqFA9Sg",
				IsCalculationTarget: true,
				RecordedOn:          time.Date(2020, 7, 25, 0, 0, 0, 0, time.UTC),
				Title:               "西友",
				Amount:              -3612,
				SourceName:          "セゾンカード",
				Category1Name:       "食費",
				Category2Name:       "食料品",
				Memo:                "",
				IsTransfer:          false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvCSVToDomain(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvCSVToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvCSVToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
