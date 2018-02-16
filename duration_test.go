package null

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"
)

func TestDuration_Scan(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"wrong data type (int)", fields{0, false}, args{5}, true},
		{"null-scan", fields{0, false}, args{nil}, true},
		{"bad string", fields{time.Second, false}, args{"not an interval"}, true},
		{"bad []byte]", fields{time.Second, false}, args{[]byte("not an interval")}, true},
		{"missing hour", fields{time.Second, false}, args{[]byte("what:00:00")}, true},
		{"missing minutes", fields{time.Second, false}, args{[]byte("00:what:00")}, true},
		{"missing seconds", fields{time.Second, false}, args{[]byte("00:00:what")}, true},
		{"1 second string", fields{time.Second, false}, args{"00:00:01"}, false},
		{"1 second []byte", fields{time.Second, false}, args{[]byte("00:00:01")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			if err := d.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Duration.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDuration_Value(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		{"not valid", fields{0, false}, nil, false},
		{"valid", fields{0, true}, time.Duration(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			got, err := d.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Duration.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDuration(t *testing.T) {
	type args struct {
		d     time.Duration
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Duration
	}{
		{"invalid duration", args{0, false}, Duration{Duration: 0, Valid: false}},
		{"1 second duration", args{time.Second, false}, Duration{Duration: time.Second, Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDuration(tt.args.d, tt.args.valid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationFrom(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want Duration
	}{
		{"1 second", args{time.Second}, Duration{Duration: time.Second, Valid: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationFrom(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DurationFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationFromPtr(t *testing.T) {
	sec := time.Second

	type args struct {
		d *time.Duration
	}
	tests := []struct {
		name string
		args args
		want Duration
	}{
		{"1 second", args{&sec}, Duration{Duration: time.Second, Valid: true}},
		{"nil", args{nil}, Duration{Duration: time.Duration(0), Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationFromPtr(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DurationFromPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_ValueOrZero(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{"1 second value", fields{time.Second, true}, time.Second},
		{"invalid -> want zero", fields{time.Second, false}, time.Duration(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			if got := d.ValueOrZero(); got != tt.want {
				t.Errorf("Duration.ValueOrZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"null", fields{time.Second, false}, []byte("null"), false},
		{"zero value", fields{time.Duration(0), true}, []byte("\"0s\""), false},
		{"1 minute", fields{time.Minute, true}, []byte("\"1m0s\""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			got, err := d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Duration.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"null", fields{time.Duration(0), false}, args{[]byte("null")}, false},
		{"bad value", fields{time.Duration(0), false}, args{[]byte("herp derp")}, true},
		{"bad type", fields{time.Duration(0), false}, args{[]byte("5")}, true},
		{"bad string", fields{time.Duration(0), false}, args{[]byte("\"herp derp\"")}, true},
		{"zero value string", fields{time.Duration(0), true}, args{[]byte("\"0s\"")}, false},
		{"one second string", fields{time.Second, true}, args{[]byte("\"1s\"")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			if err := d.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Duration.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDuration_SetValid(t *testing.T) {
	type fields struct {
		Duration time.Duration
		Valid    bool
	}
	type args struct {
		v time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"set invalid to one second", fields{time.Second, false}, args{time.Second}},
		{"set valid to one second", fields{time.Second, true}, args{time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{
				Duration: tt.fields.Duration,
				Valid:    tt.fields.Valid,
			}
			d.SetValid(tt.args.v)
		})
	}
}

func TestDuration_Ptr(t *testing.T) {
	d := Duration{Duration: time.Second, Valid: true}
	type fields struct {
		Duration *Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   *time.Duration
	}{
		{"nil", fields{&Duration{Duration: time.Second, Valid: false}}, nil},
		{"nil", fields{&d}, &d.Duration},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Duration.Ptr(); got != tt.want {
				t.Errorf("Duration.Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}
