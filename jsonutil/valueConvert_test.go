package jsonutil

import (
	"reflect"
	"testing"
	"time"
)

func Test_getInt(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet int64
		wantOk  bool
	}{
		{
			name:    "Valid Int",
			args:    args{value: 23},
			wantRet: 23,
			wantOk:  true,
		},
		{
			name:    "Invalid Int",
			args:    args{value: "23"},
			wantRet: 0,
			wantOk:  false,
		},
		{
			name:    "Valid Float",
			args:    args{value: float64(23)},
			wantRet: 23,
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getInt(tt.args.value)
			if gotRet != tt.wantRet {
				t.Errorf("getInt() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getInt() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_getFloat(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet float64
		wantOk  bool
	}{
		{
			name:    "Valid Float",
			args:    args{value: float64(23)},
			wantRet: 23,
			wantOk:  true,
		},
		{
			name:    "Invalid Float",
			args:    args{value: "23"},
			wantRet: 0,
			wantOk:  false,
		},
		{
			name:    "invalid Float",
			args:    args{value: int(23)},
			wantRet: 0,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getFloat(tt.args.value)
			if gotRet != tt.wantRet {
				t.Errorf("getFloat() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getFloat() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_getString(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet string
		wantOk  bool
	}{
		{
			name:    "Valid String",
			args:    args{value: "Hola"},
			wantRet: "Hola",
			wantOk:  true,
		},
		{
			name:    "invalid String",
			args:    args{value: 23},
			wantRet: "",
			wantOk:  false,
		},
		{
			name:    "nil String",
			args:    args{value: nil},
			wantRet: "",
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getString(tt.args.value)
			if gotRet != tt.wantRet {
				t.Errorf("getString() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getString() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_getBool(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet bool
		wantOk  bool
	}{
		{
			name:    "Valid Bool",
			args:    args{value: true},
			wantRet: true,
			wantOk:  true,
		},
		{
			name:    "Valid Bool",
			args:    args{value: false},
			wantRet: false,
			wantOk:  true,
		},
		{
			name:    "inalid Bool",
			args:    args{value: 23},
			wantRet: false,
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getBool(tt.args.value)
			if gotRet != tt.wantRet {
				t.Errorf("getBool() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getBool() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func toDate(date string) time.Time {
	d, _ := time.Parse(time.RFC3339, date)
	return d
}

func Test_getDate(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet time.Time
		wantOk  bool
	}{
		{
			name:    "Valid Date",
			args:    args{value: "2018-11-24T01:10:22Z"},
			wantRet: toDate("2018-11-24T01:10:22Z"),
			wantOk:  true,
		},
		{
			name:    "invalid Date",
			args:    args{value: "2018-11-24T01:10:22"},
			wantRet: toDate("0001-01-01T00:00:00Z"),
			wantOk:  false,
		},
		{
			name:    "invalid Date",
			args:    args{value: nil},
			wantRet: toDate("0001-01-01T00:00:00Z"),
			wantOk:  false,
		},
		{
			name:    "invalid Date",
			args:    args{value: 23},
			wantRet: toDate("0001-01-01T00:00:00Z"),
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getDate(tt.args.value)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("getDate() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getDate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_getMap(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet map[string]interface{}
		wantOk  bool
	}{
		{
			name:    "invalid map",
			args:    args{value: 23},
			wantRet: nil,
			wantOk:  false,
		},
		{
			name:    "Valid map",
			args:    args{value: make(map[string]interface{}, 0)},
			wantRet: make(map[string]interface{}, 0),
			wantOk:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotOk := getMap(tt.args.value)
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("getMap() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getMap() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
