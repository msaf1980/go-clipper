package clipper

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tests := []struct {
		name    string
		layout  string
		in      string
		want    time.Time
		wantErr bool
	}{
		{
			name:   "date",
			layout: "2006-01-02",
			in:     "2016-02-15",
			want:   time.Date(2016, time.February, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "date TZ",
			layout: "2006-01-02Z07:00",
			in:     "2016-02-15+05:00",
			want:   time.Date(2016, time.February, 15, 0, 0, 0, 0, time.FixedZone("", 5*3600)),
		},
		{
			name:   "date and time",
			layout: "2006-01-02T15:04:05",
			in:     "2024-03-15T10:06:21",
			want:   time.Date(2024, time.March, 15, 10, 6, 21, 0, time.UTC),
		},
		{
			name:   "date and time.nanonsec",
			layout: "2006-01-02T15:04:05.999999999",
			in:     "2024-03-15T10:06:21.1",
			want:   time.Date(2024, time.March, 15, 10, 6, 21, 1e8, time.UTC),
		},
		{
			name:   "date and time.nanonsec TZ",
			layout: "2006-01-02T15:04:05.999999999Z07:00",
			in:     "2016-02-15T10:06:21.01+05:00",
			want:   time.Date(2016, time.February, 15, 10, 6, 21, 1e7, time.FixedZone("", 5*3600)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm time.Time
			tr := newTimeValue(time.Time{}, &tm, tt.layout)
			err := tr.Set(tt.in, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeValue.Set() error = '%v', wantErr = %v", err, tt.wantErr)
			} else if tr.String() != tt.want.Format(time.RFC3339Nano) {
				t.Errorf("timeValue.Set() = '%s', want '%s'", tr.String(), tt.want.Format(time.RFC3339Nano))
			}
			s := tr.Get().(time.Time).Format(time.RFC3339Nano)
			if s != tt.want.Format(time.RFC3339Nano) {
				t.Errorf("timeValue.Set() = '%s', want '%s'", s, tt.want.Format(time.RFC3339Nano))
			}
		})
	}
}
