package services

import (
	"testing"
)

func Test_isNotNumber(t *testing.T) {
	type args struct {
		n rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid number",
			args: args{n: '7'},
			want: false,
		},
		{
			name: "Invalid number",
			args: args{n: 'x'},
			want: true,
		},
		{
			name: "Invalid number",
			args: args{n: '&'},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotNumber(tt.args.n); got != tt.want {
				t.Errorf("isNotNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLuhnAlghoritm(t *testing.T) {
	type args struct {
		orderNumber string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "Valid order number",
			args:       args{orderNumber: "5062821234567892"},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "Valid order number",
			args:       args{orderNumber: "1234561239"},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "Valid order number",
			args:       args{orderNumber: "0018"},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "Valid order number",
			args:       args{orderNumber: "141"},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:       "Invalid order number",
			args:       args{orderNumber: "506282123456789"},
			wantResult: false,
			wantErr:    false,
		},
		{
			name:       "Invalid order number: too short",
			args:       args{orderNumber: "2"},
			wantResult: false,
			wantErr:    true,
		},
		{
			name:       "Invalid order number: incorrect symbol",
			args:       args{orderNumber: "5062a821234567892"},
			wantResult: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := LuhnAlghoritm(tt.args.orderNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("LuhnAlghoritm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("LuhnAlghoritm() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestValidateOrderNumber(t *testing.T) {
	type args struct {
		orderNumber string
		algorithm   string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
		wantErr    bool
	}{
		{
			name: "Valid order number, Luhn algorithm",
			args: args{
				orderNumber: "5062821234567892",
				algorithm:   "Luhn",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Invalid order number, Luhn algorithm",
			args: args{
				orderNumber: "506282123456789",
				algorithm:   "Luhn",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "Uncorrect order number, Luhn algorithm",
			args: args{
				orderNumber: "ABCDE",
				algorithm:   "Luhn",
			},
			wantResult: false,
			wantErr:    true,
		},
		{
			name: "Any order number, None algorithm",
			args: args{
				orderNumber: "123456789",
				algorithm:   "None",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Any order number, Unnknown algorithm",
			args: args{
				orderNumber: "123456789",
				algorithm:   "Unnknown",
			},
			wantResult: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ValidateOrderNumber(tt.args.orderNumber, tt.args.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("ValidateOrderNumber() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
