package containers

import "testing"

func TestEuclideanDistance(t *testing.T) {
	type args struct {
		v1 Vector
		v2 Vector
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:    "Test1",
			args:    args{v1: Vector{1, 2, 3}, v2: Vector{1, 2, 3}},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Test2",
			args:    args{v1: Vector{1, 2, 3}, v2: Vector{1, 2, 4}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Test3",
			args:    args{v1: Vector{1, 2, 3}, v2: Vector{1, 2}},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Test4",
			args:    args{v1: Vector{1, 1}, v2: Vector{1.5, 1.5}},
			want:    0.7071067811865476,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EuclideanDistance(tt.args.v1, tt.args.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("EuclideanDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EuclideanDistance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
