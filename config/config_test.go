package config

import (
	"bigJson/pkg"
	"reflect"
	"testing"
)

func TestParseHour(t *testing.T) {
	tests := []struct {
		name       string
		timeString string
		expected   int
	}{
		{"Morning hour", "10AM", 10},
		{"Afternoon hour", "3PM", 15},
		{"Midnight", "12AM", 0},
		{"Noon", "12PM", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pkg.ParseHour(tt.timeString); got != tt.expected {
				t.Errorf("ParseHour() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestValidateKeywords(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    []string
		wantErr bool
	}{
		{
			name:    "Multiple keywords with whitespace and mixed casing",
			arg:     "Potato, Veggie, MuShRoom, Potato",
			want:    []string{"potato", "veggie", "mushroom"},
			wantErr: false,
		},
		{
			name:    "Only one keyword",
			arg:     "Potato",
			want:    []string{"potato"},
			wantErr: false,
		},
		{
			name:    "Empty keywords string",
			arg:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Only whitespace",
			arg:     "   , , ,   ",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateKeywords(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateKeywords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateKeywords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateTimeWindow(t *testing.T) {
	tests := []struct {
		name       string
		timeWindow string
		wantErr    bool
	}{
		{
			name:       "Valid time window",
			timeWindow: "10AM-1PM",
			wantErr:    false,
		},
		{
			name:       "Invalid time window - missing dash",
			timeWindow: "10AM1PM",
			wantErr:    true,
		},
		{
			name:       "Invalid time window - missing AM/PM in start time",
			timeWindow: "10-1PM",
			wantErr:    true,
		},
		{
			name:       "Invalid time window - missing AM/PM in end time",
			timeWindow: "10AM-1",
			wantErr:    true,
		},
		{
			name:       "Invalid time window - missing AM/PM in both times",
			timeWindow: "10-1",
			wantErr:    true,
		},
		{
			name:       "Overnight deliverty test",
			timeWindow: "10PM-1PM",
			wantErr:    false,
		},
		{
			name:       "Invalid end time format",
			timeWindow: "10AM-13PM",
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateTimeWindow(test.timeWindow)
			if (err != nil) != test.wantErr {
				t.Errorf("validateTimeWindow() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
