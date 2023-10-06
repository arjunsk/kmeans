package containers

import (
	"reflect"
	"testing"
)

func TestClusters_Recenter(t *testing.T) {
	tests := []struct {
		name        string
		c           Clusters
		wantErr     bool
		wantCenters []Vector
	}{
		{
			name: "Test1",
			c: Clusters{
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}},
				},
			},
			wantErr: false,
			wantCenters: []Vector{
				{2, 2},
				{1.5, 1.5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Recenter(); (err != nil) != tt.wantErr {
				t.Errorf("Recenter() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i, cluster := range tt.c {
				if cluster.center.Compare(tt.wantCenters[i]) != 0 {
					t.Errorf("Recenter() gotCenter = %v, want %v", cluster.center, tt.wantCenters[i])
				}
			}
		})
	}
}

func TestClusters_RecenterWithDeltaDistance(t *testing.T) {
	type args struct {
		distFn DistanceFunction
	}
	tests := []struct {
		name              string
		c                 Clusters
		args              args
		wantMoveDistances []float64
		wantErr           bool
		wantCenters       []Vector
	}{
		{
			name: "Test1",
			c: Clusters{
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}},
				},
			},
			wantErr: false,
			wantMoveDistances: []float64{
				1.4142135623730951,
				0.7071067811865476,
			},
			wantCenters: []Vector{
				{2, 2},
				{1.5, 1.5},
			},
			args: args{
				distFn: EuclideanDistance,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMoveDistances, err := tt.c.RecenterWithDeltaDistance(tt.args.distFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecenterWithDeltaDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMoveDistances, tt.wantMoveDistances) {
				t.Errorf("RecenterWithDeltaDistance() gotMoveDistances = %v, want %v", gotMoveDistances, tt.wantMoveDistances)
			}
			for i, cluster := range tt.c {
				if cluster.center.Compare(tt.wantCenters[i]) != 0 {
					t.Errorf("Recenter() gotCenter = %v, want %v", cluster.center, tt.wantCenters[i])
				}
			}
		})
	}
}

func TestClusters_Reset(t *testing.T) {
	tests := []struct {
		name string
		c    Clusters
	}{
		{
			name: "Test1",
			c: Clusters{
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
				Cluster{
					center:  Vector{2, 2},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Reset()
			for _, cluster := range tt.c {
				if len(cluster.members) != 0 {
					t.Errorf("Clusters.Reset() = %v, want %v", cluster.members, []Vector{})
				}
				if cluster.center.Compare(Vector{}) == 0 {
					// If center is cleared, then there is a problem.
					t.Errorf("Clusters.Reset() = %v, want %v", Vector{}, cluster.center)
				}
			}
		})
	}
}

func TestClusters_Nearest(t *testing.T) {
	type args struct {
		point  Vector
		distFn DistanceFunction
	}
	tests := []struct {
		name              string
		c                 Clusters
		args              args
		wantMinClusterIdx int
		wantMinDistance   float64
		wantErr           bool
	}{
		{
			name: "Test1",
			c: Clusters{
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
				Cluster{
					center:  Vector{2, 2},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
			},
			args: args{
				point:  Vector{1, 1},
				distFn: EuclideanDistance,
			},
			wantMinClusterIdx: 0,
			wantMinDistance:   0,
			wantErr:           false,
		},
		{
			name: "Test2",
			c: Clusters{
				Cluster{
					center:  Vector{1, 1},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
				Cluster{
					center:  Vector{2, 2},
					members: []Vector{{1, 1}, {2, 2}, {3, 3}},
				},
			},
			args: args{
				point:  Vector{3, 3},
				distFn: EuclideanDistance,
			},
			wantMinClusterIdx: 1,
			wantMinDistance:   1.4142135623730951,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMinClusterIdx, gotMinDistance, err := tt.c.Nearest(tt.args.point, tt.args.distFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Nearest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMinClusterIdx != tt.wantMinClusterIdx {
				t.Errorf("Nearest() gotMinClusterIdx = %v, want %v", gotMinClusterIdx, tt.wantMinClusterIdx)
			}
			if gotMinDistance != tt.wantMinDistance {
				t.Errorf("Nearest() gotMinDistance = %v, want %v", gotMinDistance, tt.wantMinDistance)
			}
		})
	}
}
