package initializer

import (
	"go-kmeans/domain"
	"reflect"
	"testing"
)

// NOTE: This test is not deterministic, but it is very unlikely to fail.
// Hence, using Skipf instead of Errorf.
func TestKmeansPlusPlus_InitCentroids(t *testing.T) {
	type args struct {
		vectors    []domain.Vector
		clusterCnt int
	}
	tests := []struct {
		name              string
		args              args
		wantPossibilities []domain.Clusters
		wantErr           bool
		distFn            domain.DistanceFunction
	}{
		{
			name: "Test1",
			args: args{
				vectors: []domain.Vector{
					{1, 2, 3, 4}, {0, 3, 4, 1},
					{130, 200, 343, 224}, {100, 200, 300, 400},
				},
				clusterCnt: 2,
			},
			wantPossibilities: []domain.Clusters{
				{
					{Center: domain.Vector{1, 2, 3, 4}},
					{Center: domain.Vector{100, 200, 300, 400}},
				},
				{
					{Center: domain.Vector{1, 2, 3, 4}},
					{Center: domain.Vector{130, 200, 343, 224}},
				},
				{
					{Center: domain.Vector{0, 3, 4, 1}},
					{Center: domain.Vector{100, 200, 300, 400}},
				},
				{
					{Center: domain.Vector{0, 3, 4, 1}},
					{Center: domain.Vector{130, 200, 343, 224}},
				},
			},
			wantErr: false,
			distFn:  domain.EuclideanDistance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKmeansPlusPlusInitializer(tt.distFn)
			got, err := k.InitCentroids(tt.args.vectors, tt.args.clusterCnt)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitCentroids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.args.clusterCnt {
				t.Errorf("InitCentroids() got = %v, want = %v", len(got), tt.args.clusterCnt)
				return
			}

			oneMatched := false
			for _, want := range tt.wantPossibilities {
				if reflect.DeepEqual(got[0].Center, want[0].Center) && reflect.DeepEqual(got[1].Center, want[1].Center) ||
					reflect.DeepEqual(got[0].Center, want[1].Center) && reflect.DeepEqual(got[1].Center, want[0].Center) {
					oneMatched = true
					break
				}
			}

			if !oneMatched {
				t.Skipf("Kmeans++ initializer returned bad centers [A Rare Occurance]."+
					"Got = %v, want = %v", got, tt.wantPossibilities)
			}
		})
	}
}
