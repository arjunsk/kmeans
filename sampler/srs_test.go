package sampler

import (
	"testing"
)

func TestSrsSampling(t *testing.T) {
	type args[T any] struct {
		input   []T
		percent float64
	}
	type testCase[T any] struct {
		name      string
		args      args[T]
		wantCount int
	}
	tests := []testCase[int]{
		{
			name: "Test1",
			args: struct {
				input   []int
				percent float64
			}{input: []int{1, 2, 3, 4, 5, 6}, percent: 50},
			wantCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SrsSampling(tt.args.input, tt.args.percent); len(got) != tt.wantCount {
				t.Errorf("SrsSampling() actual wantCount = %v, expected wantCount %v", len(got), tt.wantCount)
			}
		})
	}
}
