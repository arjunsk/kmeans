package initializer

import "testing"

func Test_validateArgs(t *testing.T) {
	type args struct {
		vectors    [][]float64
		clusterCnt int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1 - nil vectors",
			args: args{
				vectors:    nil,
				clusterCnt: 1,
			},
			wantErr: true,
		},
		{
			name: "Test 2 - empty vectors",
			args: args{
				vectors:    make([][]float64, 0),
				clusterCnt: 1,
			},
			wantErr: true,
		},
		{
			name: "Test 3 - no error",
			args: args{
				vectors:    [][]float64{{1, 2, 3}, {4, 5, 6}},
				clusterCnt: 1,
			},
			wantErr: false,
		},
		{
			name: "Test 3 - k = 0",
			args: args{
				vectors:    [][]float64{{1, 2, 3}, {4, 5, 6}},
				clusterCnt: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateArgs(tt.args.vectors, tt.args.clusterCnt); (err != nil) != tt.wantErr {
				t.Errorf("validateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
