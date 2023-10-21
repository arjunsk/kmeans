package clusterer

import "testing"

func Test_validateArgs(t *testing.T) {
	type args struct {
		vectors            [][]float64
		clusterCnt         int
		deltaThreshold     float64
		iterationThreshold int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1 - dimension mismatch",
			args: args{
				vectors:            [][]float64{{1, 2, 3}, {1, 2}},
				clusterCnt:         1,
				deltaThreshold:     0.7,
				iterationThreshold: 100,
			},
			wantErr: true,
		},
		{
			name: "Test 2 - k > len(vectors)",
			args: args{
				vectors:            [][]float64{{1, 2, 3}, {1, 2, 4}},
				clusterCnt:         3,
				deltaThreshold:     0.7,
				iterationThreshold: 100,
			},
			wantErr: true,
		},
		{
			name: "Test 3 - delta threshold",
			args: args{
				vectors:            [][]float64{{1, 2, 3}, {1, 2, 4}},
				clusterCnt:         3,
				deltaThreshold:     1.1,
				iterationThreshold: 100,
			},
			wantErr: true,
		},
		{
			name: "Test 4 - iteration threshold",
			args: args{
				vectors:            [][]float64{{1, 2, 3}, {1, 2, 4}},
				clusterCnt:         3,
				deltaThreshold:     0.9,
				iterationThreshold: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateArgs(tt.args.vectors, tt.args.clusterCnt, tt.args.deltaThreshold, tt.args.iterationThreshold); (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
