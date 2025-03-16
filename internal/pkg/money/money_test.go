package money

import "testing"

func TestCalculatePercentageVariation(t *testing.T) {
	type args struct {
		curr int64
		prev int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Test 1",
			args: args{
				curr: 1000,
				prev: 500,
			},
			want: 100_00,
		},
		{
			name: "Test 2",
			args: args{
				curr: 500,
				prev: 1000,
			},
			want: -50_00,
		},
		{
			name: "Test 3",
			args: args{
				curr: 500,
				prev: 500,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePercentageVariation(tt.args.curr, tt.args.prev); got != tt.want {
				t.Errorf(
					"CalculatePercentageVariation() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
