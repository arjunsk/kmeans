package initializer

import (
	"github.com/arjunsk/kmeans/containers"
	"reflect"
	"testing"
)

// NOTE: This test is not used to check the output, but to compare the center choices with Kmeans++ initializer.
// Hence, using Logf instead of Errorf.
// NOTE: Kmeans initializer will return bad centers and will often throw warning more than Kmeans++ initializer test.
// This is expected behaviour. You can try re-running the test at package level, to see the difference.
func TestInitCentroids_random(t *testing.T) {
	tests := []IO{
		genIO1(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewRandomInitializer()
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
				if reflect.DeepEqual(got[0].Center(), want[0].Center()) && reflect.DeepEqual(got[1].Center(), want[1].Center()) ||
					reflect.DeepEqual(got[0].Center(), want[1].Center()) && reflect.DeepEqual(got[1].Center(), want[0].Center()) {
					oneMatched = true
					break
				}
			}

			if !oneMatched {
				t.Logf("Kmeans initializer returned bad centers [Expected Behaviour]."+
					"Got = %v, want = %v", got, tt.wantPossibilities)
			}
		})
	}
}

// NOTE: This test is not deterministic, but it is very unlikely to fail.
// Hence, using Logf instead of Errorf.
func TestInitCentroids_kmeansPlusPlus(t *testing.T) {

	tests := []struct {
		IO
		distFn containers.DistanceFunction
	}{
		{
			IO:     genIO1(),
			distFn: containers.EuclideanDistance,
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
				if reflect.DeepEqual(got[0].Center(), want[0].Center()) && reflect.DeepEqual(got[1].Center(), want[1].Center()) ||
					reflect.DeepEqual(got[0].Center(), want[1].Center()) && reflect.DeepEqual(got[1].Center(), want[0].Center()) {
					oneMatched = true
					break
				}
			}

			if !oneMatched {
				t.Logf("Kmeans++ initializer returned bad centers [A Rare Occurance]."+
					"Got = %v, want = %v", got, tt.wantPossibilities)
			}
		})
	}
}

func genIO1() IO {
	return IO{
		name: "Test1",
		args: Args{
			vectors: [][]float64{
				{1, 2, 3, 4}, {0, 3, 4, 1},
				{130, 200, 343, 224}, {100, 200, 300, 400},
			},
			clusterCnt: 2,
		},
		wantPossibilities: []containers.Clusters{
			{
				containers.NewCluster(containers.Vector{1, 2, 3, 4}),
				containers.NewCluster(containers.Vector{100, 200, 300, 400}),
			},
			{
				containers.NewCluster(containers.Vector{1, 2, 3, 4}),
				containers.NewCluster(containers.Vector{130, 200, 343, 224}),
			},
			{
				containers.NewCluster(containers.Vector{0, 3, 4, 1}),
				containers.NewCluster(containers.Vector{100, 200, 300, 400}),
			},
			{
				containers.NewCluster(containers.Vector{0, 3, 4, 1}),
				containers.NewCluster(containers.Vector{130, 200, 343, 224}),
			},
		},
		wantErr: false,
	}
}

type Args struct {
	vectors    [][]float64
	clusterCnt int
}

type IO struct {
	name              string
	args              Args
	wantErr           bool
	wantPossibilities []containers.Clusters
}
