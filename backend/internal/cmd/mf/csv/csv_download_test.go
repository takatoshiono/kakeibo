package csv

import (
	"reflect"
	"testing"
	"time"
)

func TestNewCmdCSVDownload(t *testing.T) {
	type args struct {
		o *DownloadOption
		a []string
	}
	tests := []struct {
		name string
		args args
		want *DownloadOption
	}{
		{
			name: "ok",
			args: args{
				o: &DownloadOption{},
				a: []string{
					"--year", "2020",
					"--month", "10",
					"--from", "202010",
					"--to", "202011",
					"--filename", "test.csv",
				},
			},
			want: &DownloadOption{
				year:     2020,
				month:    10,
				from:     "202010",
				to:       "202011",
				fileName: "test.csv",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCmdCSVDownload(tt.args.o)
			if err := c.ParseFlags(tt.args.a); err != nil {
				t.Errorf("c.ParseFlags() = %v, want nil", err)
			}
			if got := tt.args.o; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadOption_Validate(t *testing.T) {
	type fields struct {
		MoneyForwardEmail    string
		MoneyForwardPassword string
		year                 int
		month                int
		from                 string
		to                   string
		fileName             string
		fromTime             time.Time
		toTime               time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				year:     1,
				month:    1,
				from:     "",
				to:       "",
				fileName: "out.csv",
			},
			wantErr: false,
		},
		{
			name: "ok if from and to are both set",
			fields: fields{
				year:     1,
				month:    1,
				from:     "202001",
				to:       "202011",
				fileName: "out.csv",
			},
			wantErr: false,
		},
		{
			name: "error if year is 0",
			fields: fields{
				year:     0,
				month:    1,
				from:     "",
				to:       "",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if month is less than 1",
			fields: fields{
				year:     1,
				month:    0,
				from:     "",
				to:       "",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if month is larger than 12",
			fields: fields{
				year:     1,
				month:    13,
				from:     "",
				to:       "",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if from is set and to is not set",
			fields: fields{
				year:     1,
				month:    1,
				from:     "202001",
				to:       "",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if from is not set and to is set",
			fields: fields{
				year:     1,
				month:    1,
				from:     "",
				to:       "202011",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if from is invalid format",
			fields: fields{
				year:     1,
				month:    1,
				from:     "2020-01",
				to:       "202011",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if to is invalid format",
			fields: fields{
				year:     1,
				month:    1,
				from:     "202001",
				to:       "2020-11",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if to is before from",
			fields: fields{
				year:     1,
				month:    1,
				from:     "202001",
				to:       "201912",
				fileName: "out.csv",
			},
			wantErr: true,
		},
		{
			name: "error if filename is not set",
			fields: fields{
				year:     1,
				month:    1,
				from:     "",
				to:       "",
				fileName: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &DownloadOption{
				MoneyForwardEmail:    tt.fields.MoneyForwardEmail,
				MoneyForwardPassword: tt.fields.MoneyForwardPassword,
				year:                 tt.fields.year,
				month:                tt.fields.month,
				from:                 tt.fields.from,
				to:                   tt.fields.to,
				fileName:             tt.fields.fileName,
				fromTime:             tt.fields.fromTime,
				toTime:               tt.fields.toTime,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadOption.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownloadOption_Parse(t *testing.T) {
	type fields struct {
		MoneyForwardEmail    string
		MoneyForwardPassword string
		year                 int
		month                int
		from                 string
		to                   string
		fileName             string
		fromTime             time.Time
		toTime               time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok if fromTime and toTime are both set",
			fields: fields{
				fromTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				toTime:   time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "ok if fromTime and toTime are not both set",
			fields: fields{
				fromTime: time.Time{},
				toTime:   time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &DownloadOption{
				MoneyForwardEmail:    tt.fields.MoneyForwardEmail,
				MoneyForwardPassword: tt.fields.MoneyForwardPassword,
				year:                 tt.fields.year,
				month:                tt.fields.month,
				from:                 tt.fields.from,
				to:                   tt.fields.to,
				fileName:             tt.fields.fileName,
				fromTime:             tt.fields.fromTime,
				toTime:               tt.fields.toTime,
			}
			if err := o.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadOption.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownloadOption_convertFileName(t *testing.T) {
	type fields struct {
		MoneyForwardEmail    string
		MoneyForwardPassword string
		year                 int
		month                int
		from                 string
		to                   string
		fileName             string
		fromTime             time.Time
		toTime               time.Time
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "file name with extension",
			fields: fields{
				fileName: "foo.txt",
			},
			args: args{
				t: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC),
			},
			want: "foo-202011.txt",
		},
		{
			name: "file name without extension",
			fields: fields{
				fileName: "foo",
			},
			args: args{
				t: time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC),
			},
			want: "foo-202011",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &DownloadOption{
				MoneyForwardEmail:    tt.fields.MoneyForwardEmail,
				MoneyForwardPassword: tt.fields.MoneyForwardPassword,
				year:                 tt.fields.year,
				month:                tt.fields.month,
				from:                 tt.fields.from,
				to:                   tt.fields.to,
				fileName:             tt.fields.fileName,
				fromTime:             tt.fields.fromTime,
				toTime:               tt.fields.toTime,
			}
			if got := o.convertFileName(tt.args.t); got != tt.want {
				t.Errorf("DownloadOption.convertFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
